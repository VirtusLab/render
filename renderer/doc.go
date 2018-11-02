/*
Package renderer implements data-driven templates for generating textual output

The renderer extends the standard golang text/template and sprig functions.

Templates are executed by applying them to a data structure (configuration).
Values in the template refer to elements of the data structure (typically a field of a struct or a key in a map).

Actions can be combined using UNIX-like pipelines.

The input text for a template is UTF-8-encoded text in any format.

See renderer.ExtraFunctions for our custom functions.

Detailed documentation on the syntax and available functions can be found here:

  * https://golang.org/pkg/text/template
  * http://masterminds.github.io/sprig
  * https://godoc.org/github.com/VirtusLab/render/renderer#Renderer.ExtraFunctions

*/
package renderer
