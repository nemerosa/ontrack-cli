package utils

import (
	"errors"
	"os"

	"github.com/spf13/cobra"
)

func GetProjectFlag(cmd *cobra.Command) (string, error) {
	project, err := cmd.Flags().GetString("project")
	if err != nil {
		return "", err
	}
	if project == "" {
		project = os.Getenv("YONTRACK_PROJECT_NAME")
	}
	if project == "" {
		return "", errors.New("project is required (use --project flag or YONTRACK_PROJECT_NAME environment variable)")
	} else {
		return project, nil
	}
}

func GetBranchFlag(cmd *cobra.Command, ignoreEmptyBranch bool, normalizeBranch bool) (string, error) {
	branch, err := cmd.Flags().GetString("branch")
	if err != nil {
		return "", err
	}
	if branch == "" {
		branch = os.Getenv("YONTRACK_BRANCH_NAME")
	}
	if branch == "" && !ignoreEmptyBranch {
		return "", errors.New("branch is required (use --branch flag or YONTRACK_BRANCH_NAME environment variable)")
	} else if normalizeBranch {
		return NormalizeBranchName(branch), nil
	} else {
		return branch, nil
	}
}

func GetProjectBranchFlags(cmd *cobra.Command, ignoreEmpty bool, normalizeBranch bool) (string, string, error) {
	project, err := GetProjectFlag(cmd)
	if err != nil {
		return "", "", err
	}
	branch, err := GetBranchFlag(cmd, ignoreEmpty, normalizeBranch)
	if err != nil {
		return "", "", err
	}
	return project, branch, nil
}

func GetBuildFlag(cmd *cobra.Command) (string, error) {
	build, err := cmd.Flags().GetString("build")
	if err != nil {
		return "", err
	}
	if build == "" {
		build = os.Getenv("YONTRACK_BUILD_NAME")
	}
	if build == "" {
		return "", errors.New("build is required (use --build flag or YONTRACK_BUILD_NAME environment variable)")
	} else {
		return build, nil
	}
}

func GetProjectBranchBuildFlags(cmd *cobra.Command, ignoreEmptyBranch bool, normalizeBranch bool) (string, string, string, error) {
	project, branch, err := GetProjectBranchFlags(cmd, ignoreEmptyBranch, normalizeBranch)
	if err != nil {
		return "", "", "", err
	}
	build, err := GetBuildFlag(cmd)
	if err != nil {
		return "", "", "", err
	}
	return project, branch, build, nil
}
