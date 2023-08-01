package installselectors

import (
	"testing"

	"github.com/Masterminds/semver/v3"
	"github.com/openshift/osde2e/pkg/common/spi"
)

func TestOldestClusterImageSetSelectVersion(t *testing.T) {
	tests := []struct {
		name            string
		versions        *spi.VersionList
		expectedVersion *semver.Version
		expectedErr     bool
	}{
		{
			name: "get oldest version",
			versions: spi.NewVersionListBuilder().
				AvailableVersions([]*spi.Version{
					spi.NewVersionBuilder().Version(semver.MustParse("4.1.0")).Build(),
					spi.NewVersionBuilder().Version(semver.MustParse("4.2.0")).Build(),
					spi.NewVersionBuilder().Default(true).Version(semver.MustParse("4.3.0")).Build(),
					spi.NewVersionBuilder().Version(semver.MustParse("4.4.0")).Build(),
					spi.NewVersionBuilder().Version(semver.MustParse("4.5.0")).Build(),
				}).
				Build(),
			expectedVersion: semver.MustParse("4.1.0"),
		},
		{
			name: "get oldest version out of order",
			versions: spi.NewVersionListBuilder().
				AvailableVersions([]*spi.Version{
					spi.NewVersionBuilder().Default(true).Version(semver.MustParse("4.3.0")).Build(),
					spi.NewVersionBuilder().Version(semver.MustParse("4.2.0")).Build(),
					spi.NewVersionBuilder().Version(semver.MustParse("4.5.0")).Build(),
					spi.NewVersionBuilder().Version(semver.MustParse("4.4.0")).Build(),
					spi.NewVersionBuilder().Version(semver.MustParse("4.1.0")).Build(),
				}).
				Build(),
			expectedVersion: semver.MustParse("4.1.0"),
		},
		{
			name: "no versions",
			versions: spi.NewVersionListBuilder().
				AvailableVersions([]*spi.Version{}).
				Build(),
			expectedVersion: nil,
			expectedErr:     true,
		},
	}

	for _, test := range tests {
		selector := oldestClusterImageSet{}
		selectedVersion, descriptor, err := selector.SelectVersion(test.versions)

		if err != nil && !test.expectedErr {
			t.Errorf("test %s: error while selecting version: %v", test.name, err)
		}

		if err == nil {
			expectedDescriptor := "oldest version"
			if descriptor != expectedDescriptor {
				t.Errorf("test %s: descriptor (%s) does not match expected '%s'", test.name, descriptor, expectedDescriptor)
			}

			if (selectedVersion == nil || test.expectedVersion == nil) && selectedVersion != test.expectedVersion {
				t.Errorf("test %s: expected selected version (%v) to match expected version (%v) and one is nil", test.name, selectedVersion, test.expectedVersion)
			} else if selectedVersion != nil && !selectedVersion.Equal(test.expectedVersion) {
				t.Errorf("test %s: selected version (%v) does not match expected version (%v)", test.name, selectedVersion, test.expectedVersion)
			}
		}
	}
}
