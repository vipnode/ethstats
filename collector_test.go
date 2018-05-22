package main

import (
	"testing"
	"time"

	"github.com/vipnode/ethstats/stats"
)

func TestCollector(t *testing.T) {
	col := collector{}
	if err := col.Collect(stats.PingReport{"foo", time.Now()}); err != ErrNodeNotAuthorized {
		t.Errorf("collected unauthorized report: err=%q", err)
	}

	if err := col.Collect(stats.AuthReport{ID: "foo"}); err != nil {
		t.Errorf("failed to collect auth: err=%q", err)
	}

	if err := col.Collect(stats.PingReport{"foo", time.Now()}); err != nil {
		t.Errorf("failed to collect ping after auth: err=%q", err)
	}

}
