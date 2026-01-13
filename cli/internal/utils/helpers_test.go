package utils

import (
	"testing"

	"github.com/wordsail/cli/pkg/models"
)

func TestFindServerByName(t *testing.T) {
	servers := []models.Server{
		{Name: "server1", IP: "1.1.1.1"},
		{Name: "server2", IP: "2.2.2.2"},
		{Name: "server3", IP: "3.3.3.3"},
	}

	tests := []struct {
		name       string
		serverName string
		wantNil    bool
		wantIP     string
	}{
		{"find first server", "server1", false, "1.1.1.1"},
		{"find middle server", "server2", false, "2.2.2.2"},
		{"find last server", "server3", false, "3.3.3.3"},
		{"not found", "server4", true, ""},
		{"empty name", "", true, ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := FindServerByName(servers, tt.serverName)
			if tt.wantNil && result != nil {
				t.Errorf("FindServerByName() = %v, want nil", result)
			}
			if !tt.wantNil && result == nil {
				t.Errorf("FindServerByName() = nil, want server with IP %s", tt.wantIP)
			}
			if !tt.wantNil && result != nil && result.IP != tt.wantIP {
				t.Errorf("FindServerByName() IP = %v, want %v", result.IP, tt.wantIP)
			}
		})
	}
}

func TestFindServerIndexByName(t *testing.T) {
	servers := []models.Server{
		{Name: "server1"},
		{Name: "server2"},
		{Name: "server3"},
	}

	tests := []struct {
		name       string
		serverName string
		wantIndex  int
	}{
		{"find first", "server1", 0},
		{"find middle", "server2", 1},
		{"find last", "server3", 2},
		{"not found", "server4", -1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := FindServerIndexByName(servers, tt.serverName)
			if result != tt.wantIndex {
				t.Errorf("FindServerIndexByName() = %v, want %v", result, tt.wantIndex)
			}
		})
	}
}

func TestFindSiteBySystemName(t *testing.T) {
	server := &models.Server{
		Name: "testserver",
		Sites: []models.Site{
			{SystemName: "site1", PrimaryDomain: "site1.com"},
			{SystemName: "site2", PrimaryDomain: "site2.com"},
		},
	}

	tests := []struct {
		name       string
		server     *models.Server
		systemName string
		wantNil    bool
		wantDomain string
	}{
		{"find first site", server, "site1", false, "site1.com"},
		{"find second site", server, "site2", false, "site2.com"},
		{"not found", server, "site3", true, ""},
		{"nil server", nil, "site1", true, ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := FindSiteBySystemName(tt.server, tt.systemName)
			if tt.wantNil && result != nil {
				t.Errorf("FindSiteBySystemName() = %v, want nil", result)
			}
			if !tt.wantNil && result == nil {
				t.Errorf("FindSiteBySystemName() = nil, want site")
			}
			if !tt.wantNil && result != nil && result.PrimaryDomain != tt.wantDomain {
				t.Errorf("FindSiteBySystemName() domain = %v, want %v", result.PrimaryDomain, tt.wantDomain)
			}
		})
	}
}

func TestFindSiteIndexBySystemName(t *testing.T) {
	server := &models.Server{
		Sites: []models.Site{
			{SystemName: "site1"},
			{SystemName: "site2"},
		},
	}

	tests := []struct {
		name       string
		server     *models.Server
		systemName string
		wantIndex  int
	}{
		{"find first", server, "site1", 0},
		{"find second", server, "site2", 1},
		{"not found", server, "site3", -1},
		{"nil server", nil, "site1", -1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := FindSiteIndexBySystemName(tt.server, tt.systemName)
			if result != tt.wantIndex {
				t.Errorf("FindSiteIndexBySystemName() = %v, want %v", result, tt.wantIndex)
			}
		})
	}
}

func TestFindSiteByDomain(t *testing.T) {
	server := &models.Server{
		Sites: []models.Site{
			{
				SystemName:    "site1",
				PrimaryDomain: "primary.com",
				Domains: []models.Domain{
					{Domain: "primary.com"},
					{Domain: "www.primary.com"},
				},
			},
			{
				SystemName:    "site2",
				PrimaryDomain: "other.com",
				Domains: []models.Domain{
					{Domain: "other.com"},
				},
			},
		},
	}

	tests := []struct {
		name           string
		server         *models.Server
		domain         string
		wantNil        bool
		wantSystemName string
	}{
		{"find by primary domain", server, "primary.com", false, "site1"},
		{"find by additional domain", server, "www.primary.com", false, "site1"},
		{"find other site", server, "other.com", false, "site2"},
		{"not found", server, "notfound.com", true, ""},
		{"nil server", nil, "primary.com", true, ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := FindSiteByDomain(tt.server, tt.domain)
			if tt.wantNil && result != nil {
				t.Errorf("FindSiteByDomain() = %v, want nil", result)
			}
			if !tt.wantNil && result == nil {
				t.Errorf("FindSiteByDomain() = nil, want site")
			}
			if !tt.wantNil && result != nil && result.SystemName != tt.wantSystemName {
				t.Errorf("FindSiteByDomain() systemName = %v, want %v", result.SystemName, tt.wantSystemName)
			}
		})
	}
}

func TestGetProvisionedServers(t *testing.T) {
	servers := []models.Server{
		{Name: "server1", Status: "provisioned"},
		{Name: "server2", Status: "unprovisioned"},
		{Name: "server3", Status: "provisioned"},
		{Name: "server4", Status: "error"},
	}

	result := GetProvisionedServers(servers)

	if len(result) != 2 {
		t.Errorf("GetProvisionedServers() returned %d servers, want 2", len(result))
	}

	for _, s := range result {
		if s.Status != "provisioned" {
			t.Errorf("GetProvisionedServers() included server with status %s", s.Status)
		}
	}
}

func TestServerExists(t *testing.T) {
	servers := []models.Server{
		{Name: "server1"},
		{Name: "server2"},
	}

	tests := []struct {
		name       string
		serverName string
		want       bool
	}{
		{"exists", "server1", true},
		{"exists 2", "server2", true},
		{"not exists", "server3", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ServerExists(servers, tt.serverName)
			if result != tt.want {
				t.Errorf("ServerExists() = %v, want %v", result, tt.want)
			}
		})
	}
}

func TestSiteExists(t *testing.T) {
	server := &models.Server{
		Sites: []models.Site{
			{SystemName: "site1"},
			{SystemName: "site2"},
		},
	}

	tests := []struct {
		name       string
		server     *models.Server
		systemName string
		want       bool
	}{
		{"exists", server, "site1", true},
		{"exists 2", server, "site2", true},
		{"not exists", server, "site3", false},
		{"nil server", nil, "site1", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := SiteExists(tt.server, tt.systemName)
			if result != tt.want {
				t.Errorf("SiteExists() = %v, want %v", result, tt.want)
			}
		})
	}
}
