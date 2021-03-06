/*
 * Copyright (c) 2019 SUSE LLC.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */

package addon

import (
	"fmt"

	"github.com/SUSE/skuba/internal/pkg/skuba/kubernetes"
	skubaconfig "github.com/SUSE/skuba/internal/pkg/skuba/skuba"
	"k8s.io/apimachinery/pkg/util/version"
)

type AddonVersionInfoUpdate struct {
	Current kubernetes.AddonsVersion
	Updated kubernetes.AddonsVersion
}

func UpdatedAddons(clusterVersion *version.Version) (AddonVersionInfoUpdate, error) {
	clusterVersionString := clusterVersion.String()
	latestAddonVersions := kubernetes.Versions[clusterVersionString].AddonsVersion

	skubaConfig, err := skubaconfig.GetSkubaConfiguration()
	if err != nil {
		return AddonVersionInfoUpdate{}, err
	}
	aviu := AddonVersionInfoUpdate{
		Current: kubernetes.AddonsVersion{},
		Updated: kubernetes.AddonsVersion{},
	}

	for addonName, version := range latestAddonVersions {
		skubaConfigVersion := skubaConfig.AddonsVersion[addonName]
		aviu.Current[addonName] = skubaConfigVersion
		if skubaConfigVersion == nil ||
			skubaConfigVersion != nil && (version.Version > skubaConfigVersion.Version || version.ManifestVersion > skubaConfigVersion.ManifestVersion) {
			aviu.Updated[addonName] = version
		}
	}
	return aviu, nil
}

func PrintAddonUpdates(updatedAddons AddonVersionInfoUpdate) {
	for addonName, versions := range updatedAddons.Updated {
		if updatedAddons.Current[addonName] == nil && updatedAddons.Updated[addonName] != nil {
			fmt.Printf("  - %s: %s (new addon)\n", addonName, versions.Version)
			continue
		}

		hasVersionUpdate := hasAddonVersionUpdateWithAddon(updatedAddons, addonName)
		hasManifestUpdate := hasAddonManifestUpdateWithAddon(updatedAddons, addonName)
		if hasVersionUpdate && !hasManifestUpdate {
			fmt.Printf("  - %s: %s -> %s\n", addonName, updatedAddons.Current[addonName].Version, versions.Version)
		} else if hasVersionUpdate || hasManifestUpdate {
			fmt.Printf("  - %s: %s -> %s (manifest version from %d to %d)\n", addonName, updatedAddons.Current[addonName].Version, versions.Version, updatedAddons.Current[addonName].ManifestVersion, versions.ManifestVersion)
		}
	}
}

func HasAddonUpdate(aviu AddonVersionInfoUpdate) bool {
	for addon, _ := range aviu.Updated {
		if hasAddonManifestUpdateWithAddon(aviu, addon) || hasAddonVersionUpdateWithAddon(aviu, addon) {
			return true
		}
	}
	return false
}

func hasAddonManifestUpdateWithAddon(aviu AddonVersionInfoUpdate, addon kubernetes.Addon) bool {
	if aviu.Current[addon] == nil {
		return true
	}
	return aviu.Updated[addon].ManifestVersion > aviu.Current[addon].ManifestVersion
}

func hasAddonVersionUpdateWithAddon(aviu AddonVersionInfoUpdate, addon kubernetes.Addon) bool {
	if aviu.Current[addon] == nil {
		return true
	}
	return aviu.Updated[addon].Version > aviu.Current[addon].Version
}
