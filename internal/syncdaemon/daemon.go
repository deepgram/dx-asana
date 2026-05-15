package syncdaemon

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/deepgram/dx-asana/internal/asana"
)

const SyncInterval = 5 * time.Minute

type Daemon struct {
	client      *asana.Client
	projectIDs  []string
	syncTicker  *time.Ticker
	done        chan bool
	lastSync    map[string]time.Time
}

func NewDaemon(apiToken string, projectIDs []string) *Daemon {
	return &Daemon{
		client:     asana.NewClient(apiToken),
		projectIDs: projectIDs,
		done:       make(chan bool),
		lastSync:   make(map[string]time.Time),
	}
}

func (d *Daemon) Start() {
	d.syncTicker = time.NewTicker(SyncInterval)
	defer d.syncTicker.Stop()

	fmt.Printf("[sync-daemon] Starting Asana sync every %v\n", SyncInterval)
	fmt.Printf("[sync-daemon] Projects: %v\n", d.projectIDs)
	fmt.Printf("[sync-daemon] Cache dir: %s\n", getCacheDir())

	// Initial sync
	d.syncAll()

	for {
		select {
		case <-d.done:
			fmt.Println("[sync-daemon] Exiting...")
			return
		case <-d.syncTicker.C:
			d.syncAll()
		}
	}
}

func (d *Daemon) Stop() {
	d.done <- true
}

func (d *Daemon) syncAll() {
	for _, projectID := range d.projectIDs {
		if err := d.syncProject(projectID); err != nil {
			fmt.Printf("[sync-daemon] Error syncing project %s: %v\n", projectID, err)
		}
	}
}

func (d *Daemon) syncProject(projectID string) error {
	tasks, err := d.client.GetTasks(projectID, asana.WithOptFields(asana.DefaultTaskFields...))
	if err != nil {
		return err
	}

	if err := cacheTasks(projectID, tasks); err != nil {
		return err
	}

	d.lastSync[projectID] = time.Now()
	fmt.Printf("[sync-daemon] ✓ Synced %d tasks for project %s\n", len(tasks), projectID)
	return nil
}

func getCacheDir() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".asana-cache")
}