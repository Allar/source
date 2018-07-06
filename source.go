package source

import (
	"fmt"
	"github.com/go-services/annotation"
	"github.com/go-services/code"
	"go/format"
)

type Source struct {
	file   *File
	parser *fileParser
}

func New(src string) (*Source, error) {
	p := newParser()
	f, err := p.parse(src)
	if err != nil {
		return nil, err
	}
	return &Source{
		file:   f,
		parser: p,
	}, nil
}

func (s *Source) AppendFieldToStruct(name string, field *code.StructField) error {
	structure, err := s.GetStructure(name)
	if err != nil {
		return err
	}
	s.file.src = appendCodeToInner(s.file.src, structure, field)
	f, err := s.parser.parse(s.file.src)
	if err != nil {
		return err
	}
	s.file = f
	return nil
}

func (s *Source) AppendMethodToInterface(name string, method *code.InterfaceMethod) error {
	inf, err := s.GetInterface(name)
	if err != nil {
		return err
	}
	s.file.src = appendCodeToInner(s.file.src, inf, method)
	f, err := s.parser.parse(s.file.src)
	if err != nil {
		return err
	}
	s.file = f
	return nil
}

func (s *Source) AppendParameterToFunction(name string, param *code.Parameter) error {
	fn, err := s.GetFunction(name)
	if err != nil {
		return err
	}
	pre := s.file.src[:fn.ParamEnd()]
	if len(fn.code.Params) > 0 {
		pre += ", "
	}
	mid := param.String()
	end := s.file.src[fn.ParamEnd():]
	s.file.src = fmt.Sprintf("%s%s%s", pre, mid, end)
	f, err := s.parser.parse(s.file.src)
	if err != nil {
		return err
	}
	s.file = f
	return nil
}

func (s *Source) AppendCodeToFunction(name string, method *code.RawCode) error {
	fn, err := s.GetFunction(name)
	if err != nil {
		return err
	}
	s.file.src = appendCodeToInner(s.file.src, fn, method)
	f, err := s.parser.parse(s.file.src)
	if err != nil {
		return err
	}
	s.file = f
	return nil
}

func (s *Source) AppendStructure(structure code.Struct) error {
	s.file.src += "\n" + structure.String()
	f, err := s.parser.parse(s.file.src)
	if err != nil {
		return err
	}
	s.file = f
	return nil
}

func (s *Source) AppendInterface(inf code.Interface) error {
	s.file.src += "\n" + inf.String()
	f, err := s.parser.parse(s.file.src)
	if err != nil {
		return err
	}
	s.file = f
	return nil
}

func (s *Source) AppendFunction(fn code.Function) error {
	s.file.src += "\n" + fn.String()
	f, err := s.parser.parse(s.file.src)
	if err != nil {
		return err
	}
	s.file = f
	return nil
}

func (s *Source) Structures() (structures []Structure) {
	for _, v := range s.file.structures {
		structures = append(structures, v)
	}
	return
}

func (s *Source) GetStructure(name string) (*Structure, error) {
	if v, ok := s.file.structures[name]; ok {
		return &v, nil
	} else {
		return nil, fmt.Errorf("no structure with name `%s` found", name)
	}
}

func (s *Source) GetInterface(name string) (*Interface, error) {
	if v, ok := s.file.interfaces[name]; ok {
		return &v, nil
	} else {
		return nil, fmt.Errorf("no interface with name `%s` found", name)
	}
}

func (s *Source) GetFunction(name string) (*Function, error) {
	if v, ok := s.file.functions[name]; ok {
		return &v, nil
	} else {
		return nil, fmt.Errorf("no function with name `%s` found", name)
	}
}

func (s *Source) Interfaces() (interfaces []Interface) {
	for _, v := range s.file.interfaces {
		interfaces = append(interfaces, v)
	}
	return
}

func (s *Source) Functions() (functions []Interface) {
	for _, v := range s.file.interfaces {
		functions = append(functions, v)
	}
	return
}

func (s *Source) AnnotateFunction(name string, ann annotation.Annotation) error {
	fn, err := s.GetFunction(name)
	if err != nil {
		return err
	}
	return s.annotate(fn, ann)
}

func (s *Source) AnnotateStructure(name string, ann annotation.Annotation) error {
	st, err := s.GetStructure(name)
	if err != nil {
		return err
	}
	return s.annotate(st, ann)
}

func (s *Source) AnnotateStructureField(structure, field string, ann annotation.Annotation) error {
	st, err := s.GetStructure(structure)
	if err != nil {
		return err
	}
	for _, f := range st.Fields() {
		if f.Name() == field {
			return s.annotate(&f, ann)
		}
	}
	return fmt.Errorf("field with name `%s` not found in structure `%s`", field, structure)
}
func (s *Source) AnnotateInterfaceMethod(inf, method string, ann annotation.Annotation) error {
	intr, err := s.GetInterface(inf)
	if err != nil {
		return err
	}
	for _, m := range intr.Methods() {
		if m.Name() == method {
			return s.annotate(&m, ann)
		}
	}
	return fmt.Errorf("method with name `%s` not found in interface `%s`", method, inf)
}

func (s *Source) AnnotateInterface(name string, ann annotation.Annotation) error {
	inf, err := s.GetInterface(name)
	if err != nil {
		return err
	}
	return s.annotate(inf, ann)
}

func (s *Source) annotate(node Node, ann annotation.Annotation) error {
	pre := s.file.src[:node.Begin()]
	mid := code.Comment(ann.String()).String() + "\n"
	end := s.file.src[node.Begin():]
	s.file.src = fmt.Sprintf(
		"%s%s%s",
		pre,
		mid,
		end,
	)
	f, err := s.parser.parse(s.file.src)
	if err != nil {
		return err
	}
	s.file = f
	return nil
}

func (s *Source) String() (string, error) {
	src, err := format.Source([]byte(s.file.src))
	return string(src), err
}