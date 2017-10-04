package environment

import (
)

func makeTerrainField(size int) TerrainField {
	logNormal := getLogNormalConverter(0.0, 0.5, 10000)
	field := TerrainField{
		SolidityCurrent: zeroSquare(size),
		SolidityConstant: applyGridDistribution(randomSquare(size), logNormal),
		SolidityAmplitude: applyGridDistribution(randomSquare(size), logNormal),
		SolidityPhase: randomSquare(size),

		HardnessCurrent: zeroSquare(size),
		HardnessConstant: applyGridDistribution(randomSquare(size), logNormal),
		HardnessAmplitude: applyGridDistribution(randomSquare(size), logNormal),
		HardnessPhase: randomSquare(size),

		MineralCapacityCurrent: zeroSquare(size),
		MineralCapacityConstant: applyGridDistribution(randomSquare(size), logNormal),
		MineralCapacityAmplitude: applyGridDistribution(randomSquare(size), logNormal),
		MineralCapacityPhase: randomSquare(size),
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
