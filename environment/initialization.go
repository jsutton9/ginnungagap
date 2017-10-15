package environment

import (
)

const ROCK_PORTION = 0.3
const BASE_HEALTH = 100.0
const BASE_CHEMICAL_CAPACITY = 100.0

func makeTerrainField(size int) TerrainField {
	logNormal := getLogNormalConverter(0.0, 0.5, 10000)
	field := TerrainField{
		SolidityCurrent: zeroSquare(size),
		SolidityConstant: applyGridDistribution(randomSquare(size), logNormal),
		SolidityAmplitude: applyGridDistribution(randomSquare(size), logNormal),
		SolidityPhase: applyGridDistribution(randomSquare(size), uniformConverter),

		HardnessCurrent: zeroSquare(size),
		HardnessConstant: applyGridDistribution(randomSquare(size), logNormal),
		HardnessAmplitude: applyGridDistribution(randomSquare(size), logNormal),
		HardnessPhase: applyGridDistribution(randomSquare(size), uniformConverter),

		ChemicalCapacityCurrent: zeroSquare(size),
		ChemicalCapacityConstant: applyGridDistribution(randomSquare(size), logNormal),
		ChemicalCapacityAmplitude: applyGridDistribution(randomSquare(size), logNormal),
		ChemicalCapacityPhase: applyGridDistribution(randomSquare(size), uniformConverter),
	}

	return field
}

func makeFluidField(size int) FluidField {
	logNormal := getLogNormalConverter(0.0, 0.5, 10000)
	field := FluidField{
		FlowBiasCurrent: zeroVectorSquare(size),
		FlowBiasConstant: applyVectorMagnitudeDistribution(
			randomVectorSquare(size), logNormal),
		FlowBiasAmplitude: applyVectorMagnitudeDistribution(
			randomVectorSquare(size), logNormal),
		FlowBiasPhase: applyVectorMagnitudeDistribution(
			randomVectorSquare(size), uniformConverter),
	}

	return field
}

func makeEnergyField(size int) EnergyField {
	logNormal := getLogNormalConverter(0.0, 0.5, 10000)
	field := EnergyField{
		PowerCurrent: zeroSquare(size),
		PowerConstant: applyGridDistribution(randomSquare(size), logNormal),
		PowerAmplitude: applyGridDistribution(randomSquare(size), logNormal),
		PowerPhase: applyGridDistribution(randomSquare(size), uniformConverter),
	}

	return field
}

func makeBlocks(size int, terrain TerrainField) [][]Block {
	solidThreshold := gridPercentile(terrain.SolidityConstant, 1.0-ROCK_PORTION)
	blocks := make([][]Block, size)
	for i:=0; i<size; i++ {
		blocks[i] = make([]Block, size)
		for j:=0; j<size; j++ {
			if terrain.SolidityConstant[i][j] > solidThreshold {
				blocks[i][j] = Block{
					Solid: true,
					ChemicalCapacity: BASE_CHEMICAL_CAPACITY*terrain.ChemicalCapacityConstant[i][j],
					Rock: rock{
						Health: BASE_HEALTH*terrain.HardnessConstant[i][j],
					},
				}
			} else {
				blocks[i][j] = Block{
					Solid: false,
					ChemicalCapacity: 1.0,
					Water: water{
						Pressure: 0.0,
					},
				}
			}
		}
	}

	// TODO: add chemicals

	return blocks
}
