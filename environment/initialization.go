package environment

import (
)

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

func makePressureField(size int) PressureField {
	logNormal := getLogNormalConverter(0.0, 0.5, 10000)
	field := PressureField{
		PressureBiasCurrent: zeroSquare(size),
		PressureBiasConstant: applyGridDistribution(randomSquare(size), logNormal),
		PressureBiasAmplitude: applyGridDistribution(randomSquare(size), logNormal),
		PressureBiasPhase: applyGridDistribution(randomSquare(size), uniformConverter),
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

func makeBlocks(size int, terrain TerrainField, constants PhysicalConstants) [][]Block {
	solidThreshold := gridPercentile(terrain.SolidityConstant, 1.0-constants.RockPortion)
	blocks := make([][]Block, size)
	for i:=0; i<size; i++ {
		blocks[i] = make([]Block, size)
		for j:=0; j<size; j++ {
			if terrain.SolidityConstant[i][j] > solidThreshold {
				blocks[i][j] = Block{
					Solid: true,
					ChemicalCapacity: constants.BaseCapacity*terrain.ChemicalCapacityConstant[i][j],
					Rock: rock{
						Health: constants.BaseHealth*terrain.HardnessConstant[i][j],
					},
				}
			} else {
				blocks[i][j] = Block{
					Solid: false,
					ChemicalCapacity: 1.0,
					Water: water{},
				}
			}
		}
	}

	logNormal := getLogNormalConverter(0.0, 0.5, 10000)
	for _, c := range []int{0, 1, 2} {
		concentrations := applyGridDistribution(randomSquare(size), logNormal)
		for i, row := range blocks {
			for j, block := range row {
				block.Chemicals[c] = block.ChemicalCapacity*concentrations[i][j]
			}
		}
	}

	return blocks
}

func makeWorld(size int, constants PhysicalConstants) World {
	world := World{
		Size: size,
		TerrainField: makeTerrainField(size),
		PressureField: makePressureField(size),
		EnergyField: makeEnergyField(size),
		Season: 0.0,
	}
	world.Blocks = makeBlocks(size, world.TerrainField, constants)
	return world
}
