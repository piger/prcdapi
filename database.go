package prcdapi

import (
	"bufio"
	"fmt"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// Moccolo is a prcd from a prcd database.
type Moccolo struct {
	// Author is the author of the prcd
	Author string

	// Text is the text of the prcd
	Text string
}

// Section is a container of Moccoli
type Section struct {
	Name    string
	Entries []Moccolo
}

type Grimoire struct {
	Sections map[string]Section
}

func NewGrimoire() *Grimoire {
	g := &Grimoire{Sections: make(map[string]Section)}
	return g
}

func (g *Grimoire) AddSection(section Section) {
	g.Sections[section.Name] = section
}

func (g *Grimoire) HasSection(name string) bool {
	_, ok := g.Sections[name]
	return ok
}

func (g *Grimoire) FromSection(name string) (*Moccolo, error) {
	if !g.HasSection(name) {
		return nil, fmt.Errorf("Section %s does not exist", name)
	}
	i := rand.Intn(len(g.Sections[name].Entries))
	return &g.Sections[name].Entries[i], nil
}

func (g *Grimoire) FromRandomSection() (*Moccolo, error) {
	sections := make([]string, 0, len(g.Sections))
	for section := range g.Sections {
		sections = append(sections, section)
	}

	i := rand.Intn(len(sections))
	section := sections[i]
	return g.FromSection(section)
}

var mre = regexp.MustCompile(`(.*?)\s*\(([^)]+)\)$`)

// LoadPrcdFile loads a text file containing a prcd database.
func LoadPrcdFile(filename string) ([]Moccolo, error) {
	result := []Moccolo{}

	file, err := os.Open(filename)
	if err != nil {
		return result, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		line = strings.Trim(line, " ")
		lline := strings.ToLower(line)
		if strings.HasPrefix(lline, "optimiz") {
			result[len(result)-1].Text = result[len(result)-1].Text + "\n" + line
		} else {
			m := mre.FindStringSubmatch(line)
			var author string
			var text string

			if m == nil || len(m) != 3 {
				author = "unknown"
				text = line
			} else {
				text = m[1]
				author = m[2]
			}

			moccolo := Moccolo{Author: author, Text: text}
			result = append(result, moccolo)
		}

	}

	if err := scanner.Err(); err != nil {
		return result, err
	}

	return result, nil
}

func LoadPrcdDir(path string) (*Grimoire, error) {
	matches, err := filepath.Glob(filepath.Join(path, "prcd_*.txt"))
	if err != nil {
		log.Print(err)
		return nil, err
	}

	grimoire := NewGrimoire()
	skips := []string{"prcd_rd!.txt", "prcd_out.txt", "prcd_vpf.txt", "prcd_hst.txt", "prcd_int.txt"}

	for _, match := range matches {
		toSkip := false

		for _, skip := range skips {
			if strings.HasSuffix(match, skip) {
				toSkip = true
				break
			}
		}
		if toSkip {
			continue
		}

		i := strings.LastIndex(match, "_")
		sectionName := match[i+1 : len(match)-4]

		entries, err := LoadPrcdFile(match)
		if err != nil {
			log.Print(err)
			return nil, err
		}

		section := Section{Name: sectionName, Entries: entries}
		grimoire.AddSection(section)
	}

	return grimoire, nil
}
