# pdf-booklet

This command line tool rearranges PDF pages for 2-sided booklet layout

Splits the provided PDF into booklets - PDFs, grouped 
by 8, 12, 16, 20, 24, 28, or 32 pages, and sorts the pages
for double side printing. On output there is a separate PDF file for each side (odd and even) of each booklet.

### Requirements

This tools requires `pdfunite` and `pdfseparate` to be installed
(from Glyph & Cog, LLC http://poppler.freedesktop.org) which typically come as a part of Ghostscript distribution

If not installed on Ubuntu this can be done `apt-get install poppler-utils`

### License
MIT
