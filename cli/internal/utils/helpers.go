package utils

import (
	"github.com/wordsail/cli/pkg/models"
)

// FindServerByName finds a server by name in the servers slice
// Returns nil if not found
func FindServerByName(servers []models.Server, name string) *models.Server {
	for i := range servers {
		if servers[i].Name == name {
			return &servers[i]
		}
	}
	return nil
}

// FindServerIndexByName finds a server's index by name in the servers slice
// Returns -1 if not found
func FindServerIndexByName(servers []models.Server, name string) int {
	for i := range servers {
		if servers[i].Name == name {
			return i
		}
	}
	return -1
}

// FindSiteBySystemName finds a site by system name within a server
// Returns nil if not found
func FindSiteBySystemName(server *models.Server, systemName string) *models.Site {
	if server == nil {
		return nil
	}
	for i := range server.Sites {
		if server.Sites[i].SystemName == systemName {
			return &server.Sites[i]
		}
	}
	return nil
}

// FindSiteIndexBySystemName finds a site's index by system name within a server
// Returns -1 if not found
func FindSiteIndexBySystemName(server *models.Server, systemName string) int {
	if server == nil {
		return -1
	}
	for i := range server.Sites {
		if server.Sites[i].SystemName == systemName {
			return i
		}
	}
	return -1
}

// FindSiteByDomain finds a site by domain (primary or additional) within a server
// Returns nil if not found
func FindSiteByDomain(server *models.Server, domain string) *models.Site {
	if server == nil {
		return nil
	}
	for i := range server.Sites {
		if server.Sites[i].PrimaryDomain == domain {
			return &server.Sites[i]
		}
		for _, d := range server.Sites[i].Domains {
			if d.Domain == domain {
				return &server.Sites[i]
			}
		}
	}
	return nil
}

// GetProvisionedServers returns only servers with status "provisioned"
func GetProvisionedServers(servers []models.Server) []models.Server {
	result := make([]models.Server, 0)
	for _, s := range servers {
		if s.Status == "provisioned" {
			result = append(result, s)
		}
	}
	return result
}

// ServerExists checks if a server with the given name exists
func ServerExists(servers []models.Server, name string) bool {
	return FindServerByName(servers, name) != nil
}

// SiteExists checks if a site with the given system name exists on a server
func SiteExists(server *models.Server, systemName string) bool {
	return FindSiteBySystemName(server, systemName) != nil
}
