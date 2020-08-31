package subscriptions

import "strings"

type paginator struct {
	builders  []*strings.Builder
	msgLength int
}

func (p *paginator) Line(s string) {
	var builder *strings.Builder

	p.msgLength = p.msgLength + len(s) + 1
	if len(p.builders) == 0 {
		builder = p.newBuilder()
	} else if p.msgLength > 1985 {
		builder = p.newBuilder()
		builder.WriteString("```")
		p.msgLength = 0

	} else {
		builder =  p.builders[len(p.builders)-1]
	}

	builder.WriteString(s)
	builder.WriteByte('\n')

}

func (p *paginator) newBuilder() *strings.Builder {
	var newBuilder strings.Builder
	newBuilder.WriteString("```\n")
	p.builders = append(p.builders, &newBuilder)
	return &newBuilder
}

func (p *paginator) GetMessages() []string {
	p.builders[len(p.builders)-1].WriteString("```")
	var messages []string
	for _, builder := range p.builders {
		messages = append(messages, builder.String())
	}

	return messages
}
