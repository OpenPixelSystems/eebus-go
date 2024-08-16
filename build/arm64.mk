################################################################################
### EEBUS Controlbox DEMO - arm64                                            ###
### Date: 16/08/2024                                                         ###
### Version: v1.0.0                                                          ###
################################################################################

BIN_ARM64 = bin/arm64
PREFIX_ARM64 = source $(ENVIRONMENT_ARM64) && CGO_ENABLED=1 GOOS=linux GOARCH=arm64

arm64: controlbox-demo-arm64

################################################################################
### EEBUS Controlbox Demo                                                    ###
################################################################################

controlbox-demo-arm64:
	@$(PREFIX_ARM64) $(GO) build $(ADD_VERSION) -o $(BIN_ARM64)/controlbox-demo openpixelsystems.org/eebus-go/cmd/controlbox/
	@echo Compiled $(BIN_ARM64)/controlbox with version \'$(COMMIT_ID)\'
