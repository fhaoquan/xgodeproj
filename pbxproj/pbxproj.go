package pbxproj

import "github.com/bitly/go-simplejson"

// Pbxproj represent project.pbxproj
type Pbxproj struct {
	path                 string
	json                 *simplejson.Json
	sections             []string
	fileReferences       []FileReference
	nativeTargets        []NativeTarget
	buildFiles           []BuildFile
	sourcesBuildPhases   []SourcesBuildPhase
	resourcesBuildPhases []ResourcesBuildPhase
	variantGroups        []VariantGroup
	groups               []Group
}

// NewPbxproj constructor
func NewPbxproj(path string) *Pbxproj {
	js := convertJSON(path)
	m := js.Get("objects").MustMap()

	return &Pbxproj{
		path,
		js,
		parseSectionNames(m),
		parseFileReferences(m),
		parseNativeTargets(m),
		parseBuildFiles(m),
		parseSourcesBuildPhases(m),
		parseResourcesBuildPhases(m),
		parseVariantGroups(m),
		parseGroups(m),
	}
}

// Exists specified section
func (p Pbxproj) Exists(section string) bool {
	return contains(p.SectionNames(), section)
}

// SectionNames return all distinct sorted section names
func (p Pbxproj) SectionNames() []string {
	return p.sections
}

// FileReferencePathNames return file reference path names
func (p Pbxproj) FileReferencePathNames() []string {
	s := []string{}
	for _, f := range p.fileReferences {
		s = append(s, f.path)
	}
	return s
}

// NativeTargetNames return all target names
func (p Pbxproj) NativeTargetNames() []string {
	s := []string{}
	for _, t := range p.nativeTargets {
		s = append(s, t.name)
	}
	return s
}

// BuildFileNames return all build file names
func (p Pbxproj) BuildFileNames() []string {
	s := []string{}
	for _, b := range p.buildFiles {
		if fr, found := p.findFileReferenceByID(b.fileRef); found {
			s = append(s, fr.path)
		}
	}
	return s
}

// BuildPhaseSourceFileNames return source file for build each target
func (p Pbxproj) BuildPhaseSourceFileNames() map[string][]string {
	m := map[string][]string{}
	for _, s := range p.sourcesBuildPhases {
		t, found := p.findNativeTargetByID(s.id)
		if !found {
			continue
		}
		m[t.name] = []string{}
		for _, id := range s.files {
			if bf, found := p.findBuildFileByID(id); found {
				if fr, found := p.findFileReferenceByID(bf.fileRef); found {
					m[t.name] = append(m[t.name], fr.path)
				}
			}
		}
	}
	return m
}

// BuildPhaseResourceFileNames return resource file for build each target
func (p Pbxproj) BuildPhaseResourceFileNames() map[string][]string {
	m := map[string][]string{}
	for _, s := range p.resourcesBuildPhases {
		t, found := p.findNativeTargetByID(s.id)
		if !found {
			continue
		}
		m[t.name] = []string{}
		for _, id := range s.files {
			if bf, found := p.findBuildFileByID(id); found {
				if fr, found := p.findFileReferenceByID(bf.fileRef); found {
					m[t.name] = append(m[t.name], fr.path)
				} else {
					if vg, found := p.findVariantGroupByID(bf.fileRef); found {
						m[t.name] = append(m[t.name], vg.name)
					}
				}
			}
		}
	}
	return m
}

// VariantGroupNames return all variant group names
func (p Pbxproj) VariantGroupNames() []string {
	s := []string{}
	for _, g := range p.variantGroups {
		s = append(s, g.name)
	}
	return s
}

// WalkFunc is the type of the function called for each fileReference or
// group visited by Walk.
type WalkFunc func(entry GroupEntry, level int)

// walk recursively descends group entry
func (p Pbxproj) walk(entry GroupEntry, level int, walkFn WalkFunc) {
	walkFn(entry, level)
	for _, c := range entry.Children(p) {
		p.walk(c, level+1, walkFn)
	}
}

// Walk walks the xcode project tree rooted at root
// calling walkFn for each group or file in the tree, including root
func (p Pbxproj) Walk(walkFn WalkFunc) {
	for _, g := range p.groups {
		if g.isRoot() {
			rootLevel := 0
			p.walk(g, rootLevel, walkFn)
		}
	}
}
