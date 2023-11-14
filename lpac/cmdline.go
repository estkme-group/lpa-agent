package lpac

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path"
	"strconv"
	"sync"
)

type CommandLine struct {
	Program string
	EnvMap  map[string]string
	mux     sync.Mutex
}

func (c *CommandLine) Info(ctx context.Context) (info *Information, err error) {
	info = new(Information)
	err = c.invoke(ctx, []string{"info"}, nil, &info)
	return
}

func (c *CommandLine) ListProfile(ctx context.Context) (profiles []*Profile, err error) {
	profiles = make([]*Profile, 0)
	err = c.invoke(ctx, []string{"profile", "list"}, nil, &profiles)
	return
}

func (c *CommandLine) SpecificProfile(ctx context.Context, iccid string) (*Profile, error) {
	profiles, err := c.ListProfile(ctx)
	if err != nil {
		return nil, err
	}
	for _, profile := range profiles {
		if profile.ICCID == iccid {
			return profile, nil
		}
	}
	return nil, nil
}

func (c *CommandLine) SetProfileName(ctx context.Context, iccid, name string) error {
	return c.invoke(ctx, []string{"profile", "rename", iccid, name}, nil, nil)
}

func (c *CommandLine) EnableProfile(ctx context.Context, iccid string) error {
	return c.invoke(ctx, []string{"profile", "enable", iccid}, nil, nil)
}

func (c *CommandLine) DisableProfile(ctx context.Context, iccid string) error {
	return c.invoke(ctx, []string{"profile", "disable", iccid}, nil, nil)
}

func (c *CommandLine) DeleteProfile(ctx context.Context, iccid string) error {
	return c.invoke(ctx, []string{"profile", "delete", iccid}, nil, nil)
}

func (c *CommandLine) ListNotification(ctx context.Context) (notifications []*Notification, err error) {
	notifications = make([]*Notification, 0)
	err = c.invoke(ctx, []string{"notification", "list"}, nil, &notifications)
	return
}

func (c *CommandLine) SpecificNotification(ctx context.Context, index int) (*Notification, error) {
	notifications, err := c.ListNotification(ctx)
	if err != nil {
		return nil, err
	}
	for _, notification := range notifications {
		if notification.Index == index {
			return notification, nil
		}
	}
	return nil, nil
}

func (c *CommandLine) ProcessNotification(ctx context.Context, index int) error {
	return c.invoke(ctx, []string{"notification", "process", strconv.Itoa(index)}, nil, nil)
}

func (c *CommandLine) RemoveNotification(ctx context.Context, index int) error {
	return c.invoke(ctx, []string{"notification", "remove", strconv.Itoa(index)}, nil, nil)
}

type DownloadProfile struct {
	SMDP        string `json:"smdp,omitempty"`
	MatchingId  string `json:"matching_id,omitempty"`
	IMEI        string `json:"imei,omitempty"`
	ConfirmCode string `json:"confirm_code,omitempty"`
}

func (c *CommandLine) DownloadProfile(ctx context.Context, cfg *DownloadProfile) error {
	var envs []string
	if cfg.SMDP != "" {
		envs = append(envs, "SMDP="+cfg.SMDP)
	}
	if cfg.MatchingId != "" {
		envs = append(envs, "MATCHINGID="+cfg.MatchingId)
	}
	if cfg.IMEI != "" {
		envs = append(envs, "IMEI="+cfg.IMEI)
	}
	if cfg.ConfirmCode != "" {
		envs = append(envs, "CONFIRMATION_CODE="+cfg.ConfirmCode)
	}
	return c.invoke(ctx, []string{"download"}, envs, nil)
}

func (c *CommandLine) SetDefaultSMDP(ctx context.Context, smdp string) error {
	return c.invoke(ctx, []string{"defaultsmdp", smdp}, nil, nil)
}

func (c *CommandLine) Purge(ctx context.Context) error {
	return c.invoke(ctx, []string{"purge", "yes"}, nil, nil)
}

func (c *CommandLine) Discovery(ctx context.Context, smdp, imei string) error {
	var envs []string
	if smdp != "" {
		envs = append(envs, "SMDP="+smdp)
	}
	if imei != "" {
		envs = append(envs, "IMEI="+imei)
	}
	return c.invoke(ctx, []string{"discovery"}, envs, nil)
}

func (c *CommandLine) invoke(ctx context.Context, args, envs []string, data any) error {
	if !c.mux.TryLock() {
		return errors.New("lpac: cannot be called concurrently")
	}
	defer c.mux.Unlock()
	cmd := exec.CommandContext(ctx, c.Program, args...)
	cmd.Dir = path.Base(cmd.Path)
	cmd.Env = os.Environ()
	cmd.Stderr = os.Stderr
	cmd.Env = append(cmd.Env, envs...)
	for name, value := range c.EnvMap {
		cmd.Env = append(cmd.Env, name+"="+value)
	}
	stdout, err := cmd.Output()
	if err != nil {
		return fmt.Errorf("lpac: %s", string(stdout))
	}
	type lpaPayload struct {
		Code    int             `json:"code"`
		Message string          `json:"message"`
		Data    json.RawMessage `json:"data"`
	}
	type lpaResponse struct {
		Type    string      `json:"type"`
		Payload *lpaPayload `json:"payload"`
	}
	var response lpaResponse
	if err = json.Unmarshal(stdout, &response); err != nil {
		return err
	}
	if response.Type != "lpa" {
		return errors.New("lpac: type error")
	}
	if p := response.Payload; p.Code == -1 {
		if p.Data == nil {
			return fmt.Errorf("lpac: %s", p.Message)
		}
		return fmt.Errorf("lpac: %s (%s)", p.Message, string(p.Data))
	}
	if data == nil {
		return nil
	}
	return json.Unmarshal(response.Payload.Data, data)
}
