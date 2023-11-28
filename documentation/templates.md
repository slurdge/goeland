Email templates
===============

Templates are provided by the standard go `"text/template"` package.
The easiest way to create a new template is to copy the original one in the `cmd/asset` subdirectory and change it according to your needs.

The following variables can be used inside the template:

|Variable|Usage|
|--------|-----|
|EntryTitle| The title of the entry. If you have a digest, it will be "Digest for" |
|EntryContent | The main body of the entry |
|EntryURL | The URL of the entry |
|IncludeHeader | Include the header logo |
|IncludeTitle | Include the header title |
|IncludeFooter| Include the footer (which can be overridden by config)|
|EntryFooter | The content of the footer |
|ContentID | The `cid:` content ID of the attachement for logo|

Create a whole HTML document. The document will be automatically converted to a text file and put as a separate attachment by goeland for text-based readers.

If you think your template is great looking and worth to be included by default, don't hesitate to make a PR!
