package reflect

// AnyMarshaller attempts to marshall as-is, and then from string
func AnyMarshaller(source interface{}, target interface{}) error {
	err := AsIsMarshaller(source, target)
	if err != nil {
		if sv, ok := source.(string); ok {
			err = FromStringMarshaller(sv, target)
		}
	}

	return err
}
