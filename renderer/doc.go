/*
Package renderer implements data-driven templates for generating textual output

The renderer extends the standard golang text/template and sprig functions.

Templates are executed by applying them to a data structure (configuration).
Values in the template refer to elements of the data structure (typically a field of a struct or a key in a map).

Actions can be combine using UNIX-like pipelines.

The input text for a template is UTF-8-encoded text in any format.

Detailed documentation of the standard functions can be found here:
- https://golang.org/pkg/text/template
- http://masterminds.github.io/sprig

See renderer.ExtraFunctions for our custom functions.

*/
package renderer
