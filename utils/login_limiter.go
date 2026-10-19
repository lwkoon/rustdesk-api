package utils

import (
	"errors"
	"sync"
	"time"
)

// security policy config
type SecurityPolicy struct {
	CaptchaThreshold int // failed attempts that trigger a captcha; below 0 means disabled, 0 means always required
	BanThreshold     int // failed attempts that trigger a ban; 0 means disabled
	AttemptsWindow   time.Duration
	BanDuration      time.Duration
}

// captcha provider interface
type CaptchaProvider interface {
	Generate() (id string, content string, answer string, err error)
	//Validate(ip, code string) bool
	Expiration() time.Duration           // captcha expiry time; should be shorter than AttemptsWindow
	Draw(content string) (string, error) // draw the captcha
}

// captcha metadata
type CaptchaMeta struct {
	Id        string
	Content   string
	Answer    string
	ExpiresAt time.Time
}

// IP ban records
type BanRecord struct {
	ExpiresAt time.Time
	Reason    string
}

// login limiter
type LoginLimiter struct {
	mu          sync.Mutex
	policy      SecurityPolicy
	attempts    map[string][]time.Time //
	captchas    map[string]CaptchaMeta
	bannedIPs   map[string]BanRecord
	provider    CaptchaProvider
	cleanupStop chan struct{}
}

var defaultSecurityPolicy = SecurityPolicy{
	CaptchaThreshold: 3,
	BanThreshold:     5,
	AttemptsWindow:   5 * time.Minute,
	BanDuration:      30 * time.Minute,
}

func NewLoginLimiter(policy SecurityPolicy) *LoginLimiter {
	// set the default values
	if policy.AttemptsWindow == 0 {
		policy.AttemptsWindow = 5 * time.Minute
	}
	if policy.BanDuration == 0 {
		policy.BanDuration = 30 * time.Minute
	}

	ll := &LoginLimiter{
		policy:      policy,
		attempts:    make(map[string][]time.Time),
		captchas:    make(map[string]CaptchaMeta),
		bannedIPs:   make(map[string]BanRecord),
		cleanupStop: make(chan struct{}),
	}
	go ll.cleanupRoutine()
	return ll
}

// register the captcha provider
func (ll *LoginLimiter) RegisterProvider(p CaptchaProvider) {
	ll.mu.Lock()
	defer ll.mu.Unlock()
	ll.provider = p
}

// isDisabled check whether the login limiter is disabled
func (ll *LoginLimiter) isDisabled() bool {
	return ll.policy.CaptchaThreshold < 0 && ll.policy.BanThreshold == 0
}

// record a failed login attempt
func (ll *LoginLimiter) RecordFailedAttempt(ip string) {
	if ll.isDisabled() {
		return
	}
	ll.mu.Lock()
	defer ll.mu.Unlock()

	if banned, _ := ll.isBanned(ip); banned {
		return
	}

	now := time.Now()
	windowStart := now.Add(-ll.policy.AttemptsWindow)

	// clear expired attempts
	validAttempts := ll.pruneAttempts(ip, windowStart)

	// record a new attempt
	validAttempts = append(validAttempts, now)
	ll.attempts[ip] = validAttempts

	// check the ban conditions
	if ll.policy.BanThreshold > 0 && len(validAttempts) >= ll.policy.BanThreshold {
		ll.banIP(ip, "excessive failed attempts")
		return
	}

	return
}

// generate a captcha
func (ll *LoginLimiter) RequireCaptcha() (error, CaptchaMeta) {
	ll.mu.Lock()
	defer ll.mu.Unlock()

	if ll.provider == nil {
		return errors.New("no captcha provider available"), CaptchaMeta{}
	}

	id, content, answer, err := ll.provider.Generate()
	if err != nil {
		return err, CaptchaMeta{}
	}

	// store the captcha
	ll.captchas[id] = CaptchaMeta{
		Id:        id,
		Content:   content,
		Answer:    answer,
		ExpiresAt: time.Now().Add(ll.provider.Expiration()),
	}

	return nil, ll.captchas[id]
}

// verify the captcha
func (ll *LoginLimiter) VerifyCaptcha(id, answer string) bool {
	ll.mu.Lock()
	defer ll.mu.Unlock()

	// find the matching captcha
	if ll.provider == nil {
		return false
	}

	// get and verify the captcha
	captcha, exists := ll.captchas[id]
	if !exists {
		return false
	}

	// clear expired captchas
	if time.Now().After(captcha.ExpiresAt) {
		delete(ll.captchas, id)
		return false
	}

	// verify and clear the state
	if answer == captcha.Answer {
		delete(ll.captchas, id)
		return true
	}

	return false
}

func (ll *LoginLimiter) DrawCaptcha(content string) (err error, str string) {
	str, err = ll.provider.Draw(content)
	return
}

// clear the record window
func (ll *LoginLimiter) RemoveAttempts(ip string) {
	ll.mu.Lock()
	defer ll.mu.Unlock()

	_, exists := ll.attempts[ip]
	if exists {
		delete(ll.attempts, ip)
	}
}

// CheckSecurityStatus check the security status
func (ll *LoginLimiter) CheckSecurityStatus(ip string) (banned bool, captchaRequired bool) {
	if ll.isDisabled() {
		return
	}
	ll.mu.Lock()
	defer ll.mu.Unlock()

	// check the ban status
	if banned, _ = ll.isBanned(ip); banned {
		return
	}

	// clear expired data
	ll.pruneAttempts(ip, time.Now().Add(-ll.policy.AttemptsWindow))

	// check the captcha requirement
	captchaRequired = len(ll.attempts[ip]) >= ll.policy.CaptchaThreshold

	return
}

// background cleanup task
func (ll *LoginLimiter) cleanupRoutine() {
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			ll.cleanupExpired()
		case <-ll.cleanupStop:
			return
		}
	}
}

// internal helper methods
func (ll *LoginLimiter) isBanned(ip string) (bool, BanRecord) {
	record, exists := ll.bannedIPs[ip]
	if !exists {
		return false, BanRecord{}
	}
	if time.Now().After(record.ExpiresAt) {
		delete(ll.bannedIPs, ip)
		return false, BanRecord{}
	}
	return true, record
}

func (ll *LoginLimiter) banIP(ip, reason string) {
	ll.bannedIPs[ip] = BanRecord{
		ExpiresAt: time.Now().Add(ll.policy.BanDuration),
		Reason:    reason,
	}
	delete(ll.attempts, ip)
	delete(ll.captchas, ip)
}

func (ll *LoginLimiter) pruneAttempts(ip string, cutoff time.Time) []time.Time {
	var valid []time.Time
	for _, t := range ll.attempts[ip] {
		if t.After(cutoff) {
			valid = append(valid, t)
		}
	}
	if len(valid) == 0 {
		delete(ll.attempts, ip)
	} else {
		ll.attempts[ip] = valid
	}
	return valid
}

func (ll *LoginLimiter) pruneCaptchas(id string) {
	if captcha, exists := ll.captchas[id]; exists {
		if time.Now().After(captcha.ExpiresAt) {
			delete(ll.captchas, id)
		}
	}
}

func (ll *LoginLimiter) cleanupExpired() {
	ll.mu.Lock()
	defer ll.mu.Unlock()

	now := time.Now()

	// clear the ban records
	for ip, record := range ll.bannedIPs {
		if now.After(record.ExpiresAt) {
			delete(ll.bannedIPs, ip)
		}
	}

	// clear the attempt records
	for ip := range ll.attempts {
		ll.pruneAttempts(ip, now.Add(-ll.policy.AttemptsWindow))
	}

	// clear the captcha
	for id := range ll.captchas {
		ll.pruneCaptchas(id)
	}
}
