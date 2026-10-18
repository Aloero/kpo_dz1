package main

import (
	"bytes"
	"strings"
	"testing"

	"vihkat/fleet"
)

type inspectionStub struct {
	decision bool
	calls    int
	seen     fleet.Transport
}

func (s *inspectionStub) Inspect(transport fleet.Transport) bool {
	s.calls++
	s.seen = transport
	return s.decision
}

func TestRunUsesInjectedInspection(t *testing.T) {
	cases := []struct {
		name        string
		operational string
		decision    bool
		want        []string
	}{
		{"accept_faulty", "0", true, []string{"Транспорт: 1", "S-1 | простота 6/10", "Самокат | инв. № S-1\n"}},
		{"reject_healthy", "1", false, []string{"Транспорт: 0", "Подходящего транспорта пока нет.", "На балансе пока нет объектов."}},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			inspection := &inspectionStub{decision: tc.decision}
			input := strings.NewReader("1\n1\nСамокат\nS-1\n6\n" + tc.operational + "\n1.5\n3\n4\n5\n0\n")
			var output bytes.Buffer

			err := run(input, &output, inspection)

			if err != nil {
				t.Fatal(err)
			}
			if inspection.calls != 1 || inspection.seen.InventoryNumber() != "S-1" {
				t.Fatal("inspection was not called exactly once for the submitted transport")
			}
			for _, want := range tc.want {
				if !strings.Contains(output.String(), want) {
					t.Fatalf("output is missing %q: %s", want, output.String())
				}
			}
		})
	}
}

func TestContainerReusesInstance(t *testing.T) {
	container := newServiceContainer()
	calls := 0
	container.register("value", func(_ *serviceContainer) interface{} {
		calls++
		return new(int)
	})
	first := container.resolve("value")
	second := container.resolve("value")
	if first != second || calls != 1 {
		t.Fatal("container did not reuse the created instance")
	}
}
