package test_vectors

var (
	// Preloaded test vector structs.
	JacobiAdditionVectors       []*JacobiAdditionVector
	JacobiDoublingVectors       []*JacobiDoublingVector
	JacobiPointVectors          []*JacobiPointVector
	AffineMultiplicationVectors []*AffineMultiplicationVector
	NegatedPointVectors         []*NegatedPointVector
	ECDSAVectors                []*ECDSAVector
)

func init() {
	var err error

	JacobiAdditionVectors, err = loadJacobiAdditionVectors()
	if err != nil {
		panic(err)
	}

	JacobiDoublingVectors, err = loadJacobiDoublingVectors()
	if err != nil {
		panic(err)
	}

	JacobiPointVectors, err = loadJacobiPointVectors()
	if err != nil {
		panic(err)
	}

	AffineMultiplicationVectors, err = loadAffineMultiplicationVectors()
	if err != nil {
		panic(err)
	}

	NegatedPointVectors, err = loadNegatedPointVectors()
	if err != nil {
		panic(err)
	}

	ECDSAVectors, err = loadECDSAVectors()
	if err != nil {
		panic(err)
	}
}
