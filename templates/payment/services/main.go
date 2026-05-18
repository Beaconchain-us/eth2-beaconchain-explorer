import (
	"context"
	"github.com/gobitfly/eth2-beaconchain-explorer/bootstrap"
	// ... other imports
)

func main() {
	// ... existing code (DB, config, etc.)

	// -------- ONLY THESE TWO LINES --------
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	if err := bootstrap.AutoInit(ctx, router); err != nil {
		logrus.Fatalf("AutoInit failed: %v", err)
	}
	// -------------------------------------

	// ... continue (start server etc.)
}