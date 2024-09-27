package redis

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"sync"
	"testing"
	"time"
	"write/common/go-session/sess"

	"github.com/stretchr/testify/assert"
)

func TestRedis(t *testing.T) {
	globalSession, err := setupSessionManager(t)
	if err != nil {
		t.Fatal(err)
	}

	go globalSession.GC()

	ctx := context.Background()

	r, _ := http.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()

	sess, err := globalSession.SessionStart(w, r)
	if err != nil {
		t.Fatal("session start failed:", err)
	}
	defer sess.SessionRelease(ctx, w)

	// SET AND GET
	err = sess.Set(ctx, "username", "astaxie")
	if err != nil {
		t.Fatal("set username failed:", err)
	}
	username := sess.Get(ctx, "username")
	if username != "astaxie" {
		t.Fatal("get username failed")
	}

	// DELETE
	err = sess.Delete(ctx, "username")
	if err != nil {
		t.Fatal("delete username failed:", err)
	}
	username = sess.Get(ctx, "username")
	if username != nil {
		t.Fatal("delete username failed")
	}

	// FLUSH
	err = sess.Set(ctx, "username", "astaxie")
	if err != nil {
		t.Fatal("set failed:", err)
	}
	err = sess.Set(ctx, "password", "1qaz2wsx")
	if err != nil {
		t.Fatal("set failed:", err)
	}
	username = sess.Get(ctx, "username")
	if username != "astaxie" {
		t.Fatal("get username failed")
	}
	password := sess.Get(ctx, "password")
	if password != "1qaz2wsx" {
		t.Fatal("get password failed")
	}
	err = sess.Flush(ctx)
	if err != nil {
		t.Fatal("flush failed:", err)
	}
	username = sess.Get(ctx, "username")
	if username != nil {
		t.Fatal("flush failed")
	}
	password = sess.Get(ctx, "password")
	if password != nil {
		t.Fatal("flush failed")
	}

	sess.SessionRelease(ctx, w)
}

func TestProvider_SessionInit(t *testing.T) {
	savePath := `
{ "save_path": "my save path", "idle_timeout": "3s"}
`
	cp := &Provider{}
	cp.SessionInit(context.Background(), 12, savePath)
	assert.Equal(t, "my save path", cp.SavePath)
	assert.Equal(t, 3*time.Second, cp.idleTimeout)
	assert.Equal(t, int64(12), cp.maxlifetime)
}

func TestStoreSessionReleaseIfPresentAndSessionDestroy(t *testing.T) {
	globalSessions, err := setupSessionManager(t)
	if err != nil {
		t.Fatal(err)
	}
	// todo test if e==nil
	go globalSessions.GC()

	ctx := context.Background()

	r, _ := http.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()

	sess, err := globalSessions.SessionStart(w, r)
	if err != nil {
		t.Fatal("session start failed:", err)
	}

	if err := globalSessions.GetProvider().SessionDestroy(ctx, sess.SessionID(ctx)); err != nil {
		t.Error(err)
		return
	}
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		sess.SessionReleaseIfPresent(ctx, httptest.NewRecorder())
	}()
	wg.Wait()
	exist, err := globalSessions.GetProvider().SessionExist(ctx, sess.SessionID(ctx))
	if err != nil {
		t.Error(err)
	}
	if exist {
		t.Fatalf("session %s should exist", sess.SessionID(ctx))
	}
}

func setupSessionManager(t *testing.T) (*sess.Manager, error) {
	redisAddr := os.Getenv("REDIS_ADDR")
	if redisAddr == "" {
		redisAddr = "127.0.0.1:6379"
	}
	redisConfig := fmt.Sprintf("%s,100,,0,30", redisAddr)

	sessionConfig := sess.NewManagerConfig(
		sess.CfgCookieName(`gosessionid`),
		sess.CfgSetCookie(true),
		sess.CfgGcLifeTime(3600),
		sess.CfgMaxLifeTime(3600),
		sess.CfgSecure(false),
		sess.CfgCookieLifeTime(3600),
		sess.CfgProviderConfig(redisConfig),
	)
	globalSessions, err := sess.NewManager("redis", sessionConfig)
	if err != nil {
		t.Log("could not create manager: ", err)
		return nil, err
	}
	return globalSessions, nil
}