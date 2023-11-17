package lpac

import (
	"bytes"
	"context"
	"encoding/json"
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

func (c *CommandLine) Info(ctx context.Context) (*Information, error) {
	data, err := c.invoke(ctx, []string{"info"}, nil)
	if err != nil {
		return nil, err
	}
	var info Information
	if err = json.Unmarshal(data, &info); err != nil {
		return nil, err
	}
	return &info, nil
}

func (c *CommandLine) ListProfile(ctx context.Context) ([]*Profile, error) {
	data, err := c.invoke(ctx, []string{"profile", "list"}, nil)
	if err != nil {
		return nil, err
	}
	var profiles []*Profile
	if err = json.Unmarshal(data, &profiles); err != nil {
		return nil, err
	}
	return profiles, nil
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
	_, err := c.invoke(ctx, []string{"profile", "rename", iccid, name}, nil)
	return err
}

func (c *CommandLine) EnableProfile(ctx context.Context, iccid string) error {
	_, err := c.invoke(ctx, []string{"profile", "enable", iccid}, nil)
	return err
}

func (c *CommandLine) DisableProfile(ctx context.Context, iccid string) error {
	_, err := c.invoke(ctx, []string{"profile", "disable", iccid}, nil)
	return err
}

func (c *CommandLine) DeleteProfile(ctx context.Context, iccid string) error {
	notifications, err := c.ListNotification(ctx)
	if err != nil {
		return err
	}
	_, err = c.invoke(ctx, []string{"profile", "delete", iccid}, nil)
	if err != nil {
		return err
	}
	for _, notification := range notifications {
		if err = c.ProcessNotification(ctx, notification.Index); err != nil {
			continue
		}
		if err = c.RemoveNotification(ctx, notification.Index); err != nil {
			continue
		}
	}
	return nil
}

func (c *CommandLine) ListNotification(ctx context.Context) ([]*Notification, error) {
	data, err := c.invoke(ctx, []string{"notification", "list"}, nil)
	if err != nil {
		return nil, err
	}
	var notifications []*Notification
	if err = json.Unmarshal(data, &notifications); err != nil {
		return nil, err
	}
	return notifications, nil
}

func (c *CommandLine) ProcessNotification(ctx context.Context, index int) error {
	_, err := c.invoke(ctx, []string{"notification", "process", strconv.Itoa(index)}, nil)
	return err
}

func (c *CommandLine) RemoveNotification(ctx context.Context, index int) error {
	_, err := c.invoke(ctx, []string{"notification", "remove", strconv.Itoa(index)}, nil)
	return err
}

type DownloadProfile struct {
	SMDP        string `json:"smdp,omitempty"`
	MatchingId  string `json:"matching_id,omitempty"`
	IMEI        string `json:"imei,omitempty"`
	ConfirmCode string `json:"confirm_code,omitempty"`
}

func (c *CommandLine) DownloadProfile(ctx context.Context, cfg *DownloadProfile) error {
	notifications, err := c.ListNotification(ctx)
	if err != nil {
		return err
	}
	_, err = c.invoke(ctx, []string{"download"}, map[string]string{
		"SMDP":              cfg.SMDP,
		"MATCHINGID":        cfg.MatchingId,
		"IMEI":              cfg.IMEI,
		"CONFIRMATION_CODE": cfg.ConfirmCode,
	})
	if err != nil {
		return err
	}
	for _, notification := range notifications {
		if err = c.ProcessNotification(ctx, notification.Index); err != nil {
			continue
		}
		if err = c.RemoveNotification(ctx, notification.Index); err != nil {
			continue
		}
	}
	return nil
}

func (c *CommandLine) SetDefaultSMDP(ctx context.Context, smdp string) error {
	_, err := c.invoke(ctx, []string{"defaultsmdp", smdp}, nil)
	return err
}

func (c *CommandLine) Purge(ctx context.Context) error {
	_, err := c.invoke(ctx, []string{"purge", "yes"}, nil)
	return err
}

func (c *CommandLine) Discovery(ctx context.Context, smdp, imei string) error {
	_, err := c.invoke(ctx, []string{"discovery"}, map[string]string{
		"SMDP": smdp,
		"IMEI": imei,
	})
	return err
}

func (c *CommandLine) invoke(ctx context.Context, args []string, envs map[string]string) (json.RawMessage, error) {
	c.mux.Lock()
	defer c.mux.Unlock()
	var stdout bytes.Buffer
	cmd := exec.CommandContext(ctx, c.Program, args...)
	if path.IsAbs(cmd.Path) {
		cmd.Dir = path.Dir(cmd.Path)
	}
	cmd.Env = os.Environ()
	cmd.Stderr = os.Stderr
	cmd.Stdout = &stdout
	for name, value := range envs {
		if value == "" {
			continue
		}
		cmd.Env = append(cmd.Env, name+"="+value)
	}
	for name, value := range c.EnvMap {
		cmd.Env = append(cmd.Env, name+"="+value)
	}
	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("lpac: %w", err)
	}
	var resp lpaResponse
	if err := json.NewDecoder(&stdout).Decode(&resp); err != nil {
		return nil, err
	}
	if p := resp.Payload; p.Code == -1 {
		return nil, fmt.Errorf("lpac: %w", &Error{Message: p.Message, Details: p.Data})
	}
	return resp.Payload.Data, nil
}
