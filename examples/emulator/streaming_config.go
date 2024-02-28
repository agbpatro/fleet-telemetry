package emulator

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"os"

	"github.com/teslamotors/fleet-telemetry/protos"
)

var (
	fieldValues = map[protos.Field][]string{
		protos.Field_DriveRail:                          []string{"something1", "something2"},
		protos.Field_ChargeState:                        []string{"something1", "something2"},
		protos.Field_BmsFullchargecomplete:              []string{"something1", "something2"},
		protos.Field_VehicleSpeed:                       []string{"something1", "something2"},
		protos.Field_Odometer:                           []string{"something1", "something2"},
		protos.Field_PackVoltage:                        []string{"something1", "something2"},
		protos.Field_PackCurrent:                        []string{"something1", "something2"},
		protos.Field_Soc:                                []string{"something1", "something2"},
		protos.Field_DCDCEnable:                         []string{"something1", "something2"},
		protos.Field_Gear:                               []string{"something1", "something2"},
		protos.Field_IsolationResistance:                []string{"something1", "something2"},
		protos.Field_PedalPosition:                      []string{"something1", "something2"},
		protos.Field_BrakePedal:                         []string{"something1", "something2"},
		protos.Field_DiStateR:                           []string{"something1", "something2"},
		protos.Field_DiHeatsinkTR:                       []string{"something1", "something2"},
		protos.Field_DiAxleSpeedR:                       []string{"something1", "something2"},
		protos.Field_DiTorquemotor:                      []string{"something1", "something2"},
		protos.Field_DiStatorTempR:                      []string{"something1", "something2"},
		protos.Field_DiVBatR:                            []string{"something1", "something2"},
		protos.Field_DiMotorCurrentR:                    []string{"something1", "something2"},
		protos.Field_Location:                           []string{"something1", "something2"},
		protos.Field_GpsState:                           []string{"something1", "something2"},
		protos.Field_GpsHeading:                         []string{"something1", "something2"},
		protos.Field_NumBrickVoltageMax:                 []string{"something1", "something2"},
		protos.Field_BrickVoltageMax:                    []string{"something1", "something2"},
		protos.Field_NumBrickVoltageMin:                 []string{"something1", "something2"},
		protos.Field_BrickVoltageMin:                    []string{"something1", "something2"},
		protos.Field_NumModuleTempMax:                   []string{"something1", "something2"},
		protos.Field_ModuleTempMax:                      []string{"something1", "something2"},
		protos.Field_NumModuleTempMin:                   []string{"something1", "something2"},
		protos.Field_ModuleTempMin:                      []string{"something1", "something2"},
		protos.Field_RatedRange:                         []string{"something1", "something2"},
		protos.Field_Hvil:                               []string{"something1", "something2"},
		protos.Field_DCChargingEnergyIn:                 []string{"something1", "something2"},
		protos.Field_DCChargingPower:                    []string{"something1", "something2"},
		protos.Field_ACChargingEnergyIn:                 []string{"something1", "something2"},
		protos.Field_ACChargingPower:                    []string{"something1", "something2"},
		protos.Field_ChargeLimitSoc:                     []string{"something1", "something2"},
		protos.Field_FastChargerPresent:                 []string{"something1", "something2"},
		protos.Field_EstBatteryRange:                    []string{"something1", "something2"},
		protos.Field_IdealBatteryRange:                  []string{"something1", "something2"},
		protos.Field_BatteryLevel:                       []string{"something1", "something2"},
		protos.Field_TimeToFullCharge:                   []string{"something1", "something2"},
		protos.Field_ScheduledChargingStartTime:         []string{"something1", "something2"},
		protos.Field_ScheduledChargingPending:           []string{"something1", "something2"},
		protos.Field_ScheduledDepartureTime:             []string{"something1", "something2"},
		protos.Field_PreconditioningEnabled:             []string{"something1", "something2"},
		protos.Field_ScheduledChargingMode:              []string{"something1", "something2"},
		protos.Field_ChargeAmps:                         []string{"something1", "something2"},
		protos.Field_ChargeEnableRequest:                []string{"something1", "something2"},
		protos.Field_ChargerPhases:                      []string{"something1", "something2"},
		protos.Field_ChargePortColdWeatherMode:          []string{"something1", "something2"},
		protos.Field_ChargeCurrentRequest:               []string{"something1", "something2"},
		protos.Field_ChargeCurrentRequestMax:            []string{"something1", "something2"},
		protos.Field_BatteryHeaterOn:                    []string{"something1", "something2"},
		protos.Field_NotEnoughPowerToHeat:               []string{"something1", "something2"},
		protos.Field_SuperchargerSessionTripPlanner:     []string{"something1", "something2"},
		protos.Field_DoorState:                          []string{"something1", "something2"},
		protos.Field_Locked:                             []string{"something1", "something2"},
		protos.Field_FdWindow:                           []string{"something1", "something2"},
		protos.Field_FpWindow:                           []string{"something1", "something2"},
		protos.Field_RdWindow:                           []string{"something1", "something2"},
		protos.Field_RpWindow:                           []string{"something1", "something2"},
		protos.Field_VehicleName:                        []string{"something1", "something2"},
		protos.Field_SentryMode:                         []string{"something1", "something2"},
		protos.Field_SpeedLimitMode:                     []string{"something1", "something2"},
		protos.Field_CurrentLimitMph:                    []string{"something1", "something2"},
		protos.Field_Version:                            []string{"something1", "something2"},
		protos.Field_TpmsPressureFl:                     []string{"something1", "something2"},
		protos.Field_TpmsPressureFr:                     []string{"something1", "something2"},
		protos.Field_TpmsPressureRl:                     []string{"something1", "something2"},
		protos.Field_TpmsPressureRr:                     []string{"something1", "something2"},
		protos.Field_SemitruckTpmsPressureRe1L0:         []string{"something1", "something2"},
		protos.Field_SemitruckTpmsPressureRe1L1:         []string{"something1", "something2"},
		protos.Field_SemitruckTpmsPressureRe1R0:         []string{"something1", "something2"},
		protos.Field_SemitruckTpmsPressureRe1R1:         []string{"something1", "something2"},
		protos.Field_SemitruckTpmsPressureRe2L0:         []string{"something1", "something2"},
		protos.Field_SemitruckTpmsPressureRe2L1:         []string{"something1", "something2"},
		protos.Field_SemitruckTpmsPressureRe2R0:         []string{"something1", "something2"},
		protos.Field_SemitruckTpmsPressureRe2R1:         []string{"something1", "something2"},
		protos.Field_TpmsLastSeenPressureTimeFl:         []string{"something1", "something2"},
		protos.Field_TpmsLastSeenPressureTimeFr:         []string{"something1", "something2"},
		protos.Field_TpmsLastSeenPressureTimeRl:         []string{"something1", "something2"},
		protos.Field_TpmsLastSeenPressureTimeRr:         []string{"something1", "something2"},
		protos.Field_InsideTemp:                         []string{"something1", "something2"},
		protos.Field_OutsideTemp:                        []string{"something1", "something2"},
		protos.Field_SeatHeaterLeft:                     []string{"something1", "something2"},
		protos.Field_SeatHeaterRight:                    []string{"something1", "something2"},
		protos.Field_SeatHeaterRearLeft:                 []string{"something1", "something2"},
		protos.Field_SeatHeaterRearRight:                []string{"something1", "something2"},
		protos.Field_SeatHeaterRearCenter:               []string{"something1", "something2"},
		protos.Field_AutoSeatClimateLeft:                []string{"something1", "something2"},
		protos.Field_AutoSeatClimateRight:               []string{"something1", "something2"},
		protos.Field_DriverSeatBelt:                     []string{"something1", "something2"},
		protos.Field_PassengerSeatBelt:                  []string{"something1", "something2"},
		protos.Field_DriverSeatOccupied:                 []string{"something1", "something2"},
		protos.Field_SemitruckPassengerSeatFoldPosition: []string{"something1", "something2"},
		protos.Field_LateralAcceleration:                []string{"something1", "something2"},
		protos.Field_LongitudinalAcceleration:           []string{"something1", "something2"},
		protos.Field_CruiseState:                        []string{"something1", "something2"},
		protos.Field_CruiseSetSpeed:                     []string{"something1", "something2"},
		protos.Field_LifetimeEnergyUsed:                 []string{"something1", "something2"},
		protos.Field_LifetimeEnergyUsedDrive:            []string{"something1", "something2"},
		protos.Field_SemitruckTractorParkBrakeStatus:    []string{"something1", "something2"},
		protos.Field_SemitruckTrailerParkBrakeStatus:    []string{"something1", "something2"},
		protos.Field_BrakePedalPos:                      []string{"something1", "something2"},
		protos.Field_RouteLastUpdated:                   []string{"something1", "something2"},
		protos.Field_RouteLine:                          []string{"something1", "something2"},
		protos.Field_MilesToArrival:                     []string{"something1", "something2"},
		protos.Field_MinutesToArrival:                   []string{"something1", "something2"},
		protos.Field_OriginLocation:                     []string{"something1", "something2"},
		protos.Field_DestinationLocation:                []string{"something1", "something2"},
		protos.Field_CarType:                            []string{"something1", "something2"},
		protos.Field_Trim:                               []string{"something1", "something2"},
		protos.Field_ExteriorColor:                      []string{"something1", "something2"},
		protos.Field_RoofColor:                          []string{"something1", "something2"},
		protos.Field_ChargePort:                         []string{"something1", "something2"},
		protos.Field_ChargePortLatch:                    []string{"something1", "something2"},
		protos.Field_Experimental_1:                     []string{"something1", "something2"},
		protos.Field_Experimental_2:                     []string{"something1", "something2"},
		protos.Field_Experimental_3:                     []string{"something1", "something2"},
		protos.Field_Experimental_4:                     []string{"something1", "something2"},
		protos.Field_GuestModeEnabled:                   []string{"something1", "something2"},
		protos.Field_PinToDriveEnabled:                  []string{"something1", "something2"},
		protos.Field_PairedPhoneKeyAndKeyFobQty:         []string{"something1", "something2"},
		protos.Field_CruiseFollowDistance:               []string{"something1", "something2"},
		protos.Field_AutomaticBlindSpotCamera:           []string{"something1", "something2"},
		protos.Field_BlindSpotCollisionWarningChime:     []string{"something1", "something2"},
		protos.Field_SpeedLimitWarning:                  []string{"something1", "something2"},
		protos.Field_ForwardCollisionWarning:            []string{"something1", "something2"},
		protos.Field_LaneDepartureAvoidance:             []string{"something1", "something2"},
		protos.Field_EmergencyLaneDepartureAvoidance:    []string{"something1", "something2"},
		protos.Field_AutomaticEmergencyBrakingOff:       []string{"something1", "something2"},
		protos.Field_LifetimeEnergyGainedRegen:          []string{"something1", "something2"},
		protos.Field_DiStateF:                           []string{"something1", "something2"},
		protos.Field_DiStateREL:                         []string{"something1", "something2"},
		protos.Field_DiStateRER:                         []string{"something1", "something2"},
		protos.Field_DiHeatsinkTF:                       []string{"something1", "something2"},
		protos.Field_DiHeatsinkTREL:                     []string{"something1", "something2"},
		protos.Field_DiHeatsinkTRER:                     []string{"something1", "something2"},
		protos.Field_DiAxleSpeedF:                       []string{"something1", "something2"},
		protos.Field_DiAxleSpeedREL:                     []string{"something1", "something2"},
		protos.Field_DiAxleSpeedRER:                     []string{"something1", "something2"},
		protos.Field_DiSlaveTorqueCmd:                   []string{"something1", "something2"},
		protos.Field_DiTorqueActualR:                    []string{"something1", "something2"},
		protos.Field_DiTorqueActualF:                    []string{"something1", "something2"},
		protos.Field_DiTorqueActualREL:                  []string{"something1", "something2"},
		protos.Field_DiTorqueActualRER:                  []string{"something1", "something2"},
		protos.Field_DiStatorTempF:                      []string{"something1", "something2"},
		protos.Field_DiStatorTempREL:                    []string{"something1", "something2"},
		protos.Field_DiStatorTempRER:                    []string{"something1", "something2"},
		protos.Field_DiVBatF:                            []string{"something1", "something2"},
		protos.Field_DiVBatREL:                          []string{"something1", "something2"},
		protos.Field_DiVBatRER:                          []string{"something1", "something2"},
		protos.Field_DiMotorCurrentF:                    []string{"something1", "something2"},
		protos.Field_DiMotorCurrentREL:                  []string{"something1", "something2"},
		protos.Field_DiMotorCurrentRER:                  []string{"something1", "something2"},
		protos.Field_EnergyRemaining:                    []string{"something1", "something2"},
		protos.Field_ServiceMode:                        []string{"something1", "something2"},
		protos.Field_BMSState:                           []string{"something1", "something2"},
		protos.Field_GuestModeMobileAccessState:         []string{"something1", "something2"},
		protos.Field_AutopilotState:                     []string{"something1", "something2"},
		protos.Field_DestinationName:                    []string{"something1", "something2"},
	}
)

type StreamingConfig struct {
	AlertTypes   []string                 `json:"alert_types,omitempty"`
	Hostname     string                   `json:"hostname,omitempty"`
	CA           string                   `json:"ca,omitempty"`
	Expiration   int64                    `json:"exp,omitempty"`
	ConfigFields map[string]*FieldSetting `json:"fields,omitempty"`
	Fields       map[protos.Field]*FieldSetting
}

type FieldSetting struct {
	IntervalSeconds int64 `json:"interval_seconds,omitempty"`
}

func NewStreamingConfig(streamingConfigPath string) (*StreamingConfig, error) {
	streamingConfig := &StreamingConfig{}
	data, err := os.ReadFile(streamingConfigPath)
	if err != nil {
		return streamingConfig, err
	}
	err = json.Unmarshal(data, streamingConfig)
	streamingConfig.Fields = make(map[protos.Field]*FieldSetting, len(streamingConfig.ConfigFields))
	for key, value := range streamingConfig.ConfigFields {
		field_int, ok := protos.Field_value[key]
		if !ok {
			return nil, fmt.Errorf("Invalid field %s", key)
		}
		streamingConfig.Fields[protos.Field(field_int)] = value
	}
	return streamingConfig, err
}

func GetFieldValue(field protos.Field) *protos.Datum {
	sampleValues, ok := fieldValues[field]
	if !ok {
		return nil
	}
	randomIndex := rand.Intn(len(sampleValues))
	return &protos.Datum{
		Key: field,
		Value: &protos.Value{
			Value: &protos.Value_StringValue{
				StringValue: sampleValues[randomIndex],
			},
		},
	}
}
