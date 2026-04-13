// Package anomaly provides a Detector that identifies port observations which
// fall outside a configured known-good baseline.
//
// Usage:
//
//	known := []scanner.Port{
//		{Protocol: "tcp", Port: 22},
//		{Protocol: "tcp", Port: 443},
//	}
//	d := anomaly.New(known)
//
//	anoms := d.Check(currentPorts)
//	for _, a := range anoms {
//		fmt.Println(a)
//	}
//
// Ports can be added to the known-good set at runtime via Add, making the
// detector suitable for adaptive baselining workflows.
package anomaly
