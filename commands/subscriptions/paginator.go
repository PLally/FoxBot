package subscriptions

import "strings"

type paginator struct {
	builders  []strings.Builder
	msgLength int
}

func (p *paginator) Line(s string) {
	if len(p.builders) == 0 {
		var newBuilder strings.Builder
		newBuilder.WriteString("```")
		p.builders = append(p.builders, newBuilder)
	}

	p.msgLength = p.msgLength + len(s) + 1
	builder := p.builders[len(p.builders)-1]
	if p.msgLength > 1990 {
		builder.WriteString("```")
		p.msgLength = 0
		var newBuilder strings.Builder
		newBuilder.WriteString("```")
		p.builders = append(p.builders, newBuilder)
	}

	builder = p.builders[len(p.builders)-1]
	builder.WriteString(s)
	builder.WriteByte('\n')


}

func (p *paginator) GetMessages() []string {
	var messages []string
	for _, builder := range p.builders {
		messages = append(messages, builder.String())
	}

	return messages
}
