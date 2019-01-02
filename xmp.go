package main

import "github.com/beevik/etree"

// getRdfDescriptions gets a slice of the description elements within the RDF Element
func getRdfDescriptions(xmp *etree.Document) []*etree.Element {
	root := xmp.SelectElement("xmpmeta")
	if root == nil {
		return nil
	}
	rdf := root.SelectElement("RDF")
	if rdf == nil {
		return nil
	}
	return rdf.SelectElements("Description")
}

// findXMPElement gets an element of the specified name in the XMP document
func findXMPElement(xmp *etree.Document, name string) *etree.Element {
	for _, description := range getRdfDescriptions(xmp) {
		foundElement := description.SelectElement(name)
		if foundElement != nil {
			return foundElement
		}
	}
	return nil
}

// findXMPAttr gets an attr of the specified name in the XMP document
func findXMPAttr(xmp *etree.Document, name string) *etree.Attr {
	for _, description := range getRdfDescriptions(xmp) {
		foundElement := description.SelectAttr(name)
		if foundElement != nil {
			return foundElement
		}
	}
	return nil
}

// filterXMPElements filters the XMP document so that any element no in the name map
// are deleted from the document
func filterXMPElements(xmp *etree.Document, names map[string]bool) {
	for _, description := range getRdfDescriptions(xmp) {
		childrenToFilter := []*etree.Element{}
		for _, child := range description.ChildElements() {
			if !names[child.Tag] {
				childrenToFilter = append(childrenToFilter, child)
			}
		}
		for _, child := range childrenToFilter {
			description.RemoveChild(child)
		}
	}
}
