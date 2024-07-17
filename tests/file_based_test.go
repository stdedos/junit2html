package tests

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/gkampitakis/go-snaps/snaps"
	"github.com/stdedos/junit2html/pkg/utils"

	"github.com/stdedos/junit2html/pkg/cmd"
	"github.com/stdedos/junit2html/pkg/logging"

	"github.com/stretchr/testify/assert"
)

const (
	SeedEnvKey       = "JUNIT2HTML_FILE_TEST_SEED"
	NoDisableLogging = "JUNIT2HTML_NO_DISABLE_LOGGING"
	SnapshotsDir     = "__snapshots__"
)

var testRand *rand.Rand

func init() {
	seed := time.Now().UnixNano()

	if os.Getenv(SeedEnvKey) == "" {
		log.Printf("Seed: %d\n", seed)
	} else {
		var err error
		seed, err = strconv.ParseInt(os.Getenv(SeedEnvKey), 10, 64)
		if err != nil {
			panic(err)
		}
	}

	testRand = rand.New(rand.NewSource(seed))

	if os.Getenv(NoDisableLogging) == "" {
		logging.SetLevel(logging.LevelOff)
	}
}

func TestMain(m *testing.M) {
	v := m.Run()

	// After all tests have run, `go-snaps` will sort snapshots.
	snaps.Clean(m, snaps.CleanOpts{Sort: true})

	os.Exit(v)
}

func TestSnapshots(t *testing.T) {
	t.Helper()
	t.Parallel()

	testDirectories, err := os.ReadDir(".")
	assert.NoError(t, err)

	for _, entry := range testDirectories {
		switch {
		case !entry.IsDir():
			continue
		case strings.HasPrefix(entry.Name(), "__"):
			continue
		}

		wd := entry

		t.Run(wd.Name(), func(tt *testing.T) {
			snapshotsDir := "./" + wd.Name() + "/" + SnapshotsDir

			files, err := inputAsGlobOrLiterally(wd)
			assert.NoError(tt, err)

			var html string

			stdoutStr, stderrStr, err := utils.CaptureOutput(func() error {
				html, err = cmd.EntryPoint(files)

				snaps.WithConfig(
					snaps.Filename("error.log"),
					snaps.Dir(snapshotsDir),
				).MatchSnapshot(tt, err)

				return nil
			})

			// We only need stderr here from now on,
			// but let's make sure we have predictable output.
			assert.Empty(tt, stdoutStr)
			assert.NoError(tt, err)

			const resultFilename = "output.html"
			snaps.WithConfig(
				snaps.Filename(resultFilename),
				snaps.Dir(snapshotsDir),
			).MatchSnapshot(tt, html)

			// Also create the "pure" HTML file
			resultDir := wd.Name() + "/result"
			err = os.MkdirAll(resultDir, 0o755)
			assert.NoError(tt, err)
			err = os.WriteFile(resultDir+"/"+resultFilename, []byte(html+"\n"), 0o644)
			assert.NoError(tt, err)

			snaps.WithConfig(
				snaps.Filename("stderr.log"),
				snaps.Dir(snapshotsDir),
			).MatchSnapshot(tt, stderrStr)
		})
	}
}

func inputAsGlobOrLiterally(dir os.DirEntry) ([]string, error) {
	var files []string

	randomFloat := testRand.Float64()

	switch {
	case randomFloat < 0.5:
		files = append(files, fmt.Sprintf("%s/*.xml", dir.Name()))
	case randomFloat >= 0.5:
		readDirList, err := os.ReadDir(dir.Name())
		if err != nil {
			return []string{}, err
		}

		for _, file := range readDirList {
			if file.IsDir() {
				continue
			}

			if !strings.HasSuffix(file.Name(), ".xml") {
				panic(fmt.Errorf("todo: remove at first inconvenience. but it is here to protect you (%s)", file.Name()))
			}

			files = append(files, fmt.Sprintf("%s/%s", dir.Name(), file.Name()))
		}
	}
	return files, nil
}
