package cli

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"

	"github.com/kevinelliott/agentmgr/pkg/config"
)

// NewConfigCommand creates the config management command group.
func NewConfigCommand(cfg *config.Config) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "config",
		Short: "Manage configuration",
		Long: `View and modify AgentManager configuration settings.

Configuration is stored in a YAML file and can be overridden with
environment variables using the AGENTMGR_ prefix.`,
	}

	cmd.AddCommand(
		newConfigShowCommand(cfg),
		newConfigGetCommand(cfg),
		newConfigSetCommand(cfg),
		newConfigPathCommand(cfg),
		newConfigInitCommand(cfg),
	)

	return cmd
}

func newConfigShowCommand(cfg *config.Config) *cobra.Command {
	return &cobra.Command{
		Use:   "show",
		Short: "Show current configuration",
		Long:  `Display the current configuration settings.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			data, err := yaml.Marshal(cfg)
			if err != nil {
				return fmt.Errorf("failed to serialize config: %w", err)
			}
			fmt.Println(string(data))
			return nil
		},
	}
}

func newConfigGetCommand(cfg *config.Config) *cobra.Command {
	return &cobra.Command{
		Use:   "get <key>",
		Short: "Get a configuration value",
		Long: `Get a configuration value by key path.

Examples:
  agentmgr config get ui.theme
  agentmgr config get updates.auto_check
  agentmgr config get logging.level`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			key := args[0]

			// Create a loader to read config
			loader := config.NewLoader()

			// Load current config
			_, err := loader.Load("")
			if err != nil {
				return fmt.Errorf("failed to load config: %w", err)
			}

			// Get the value
			value := loader.Get(key)
			if value == nil {
				return fmt.Errorf("key %q not found in configuration", key)
			}

			fmt.Printf("%s = %v\n", key, value)
			return nil
		},
	}
}

func newConfigSetCommand(cfg *config.Config) *cobra.Command {
	return &cobra.Command{
		Use:   "set <key> <value>",
		Short: "Set a configuration value",
		Long: `Set a configuration value by key path.

Examples:
  agentmgr config set ui.theme dark
  agentmgr config set updates.auto_check false
  agentmgr config set logging.level debug
  agentmgr config set ui.page_size 50
  agentmgr config set catalog.refresh_interval 2h`,
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			key := args[0]
			valueStr := args[1]

			// Create a loader to manage config
			loader := config.NewLoader()

			// Load current config (this loads the file into viper)
			_, err := loader.Load("")
			if err != nil {
				return fmt.Errorf("failed to load config: %w", err)
			}

			// Parse value based on known types for common keys
			value := parseConfigValue(key, valueStr)

			// Set the value in viper and save
			if err := loader.SetAndSave(key, value); err != nil {
				return fmt.Errorf("failed to save config: %w", err)
			}

			printSuccess("Set %s = %s", key, valueStr)
			printInfo("Config saved to %s", config.GetConfigPath())
			return nil
		},
	}
}

// parseConfigValue parses a string value into the appropriate type based on the key.
func parseConfigValue(key, value string) interface{} {
	key = strings.ToLower(key)

	// Boolean keys
	boolKeys := []string{
		"catalog.refresh_on_start",
		"updates.auto_check", "updates.notify", "updates.auto_update",
		"ui.show_hidden", "ui.use_colors", "ui.compact_mode",
		"api.enable_grpc", "api.enable_rest", "api.require_auth",
	}
	for _, k := range boolKeys {
		if key == k {
			return strings.EqualFold(value, "true") || value == "1" || strings.EqualFold(value, "yes")
		}
	}

	// Integer keys
	intKeys := []string{
		"ui.page_size",
		"api.grpc_port", "api.rest_port",
		"logging.max_size", "logging.max_age",
	}
	for _, k := range intKeys {
		if key == k {
			if i, err := strconv.Atoi(value); err == nil {
				return i
			}
			return value
		}
	}

	// Duration keys
	durationKeys := []string{
		"catalog.refresh_interval",
		"updates.check_interval",
	}
	for _, k := range durationKeys {
		if key == k {
			if d, err := time.ParseDuration(value); err == nil {
				return d
			}
			return value
		}
	}

	// Default: return as string
	return value
}

func newConfigPathCommand(cfg *config.Config) *cobra.Command {
	return &cobra.Command{
		Use:   "path",
		Short: "Show configuration file path",
		Long:  `Display the path to the configuration file.`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(config.GetConfigPath())
		},
	}
}

func newConfigInitCommand(cfg *config.Config) *cobra.Command {
	var force bool

	cmd := &cobra.Command{
		Use:   "init",
		Short: "Initialize configuration",
		Long: `Create the configuration directory and default configuration file
if they don't exist.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := config.InitConfig(); err != nil {
				return err
			}
			printSuccess("Configuration initialized at %s", config.GetConfigPath())
			return nil
		},
	}

	cmd.Flags().BoolVarP(&force, "force", "F", false, "overwrite existing config")

	return cmd
}
