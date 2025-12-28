# Goeland Filters Documentation

This document provides comprehensive documentation for all available filters in Goeland.

## Table of Contents

- [Basic Filters](#basic-filters)
- [Content Manipulation Filters](#content-manipulation-filters)
- [Specialized Filters](#specialized-filters)
- [Advanced Usage Examples](#advanced-usage-examples)

## Basic Filters

### `all`
**Description:** Default filter that includes all entries without modification.

**Usage:**
```toml
filters = ["all"]
```

**Parameters:** None

---

### `none`
**Description:** Removes all entries from the source.

**Usage:**
```toml
filters = ["none"]
```

**Parameters:** None

---

### `first`
**Description:** Keeps only the first N entries. If no number is specified, keeps only the first entry.

**Usage:**
```toml
# Keep first entry
filters = ["first"]

# Keep first 5 entries
filters = ["first(5)"]
```

**Parameters:**
- `N` (optional): Number of entries to keep (default: 1)

---

### `last`
**Description:** Keeps only the last N entries. If no number is specified, keeps only the last entry.

**Usage:**
```toml
# Keep last entry
filters = ["last"]

# Keep last 3 entries
filters = ["last(3)"]
```

**Parameters:**
- `N` (optional): Number of entries to keep (default: 1)

---

### `random`
**Description:** Keeps 1 or more random entries from the source.

**Usage:**
```toml
# Keep 1 random entry
filters = ["random"]

# Keep 5 random entries
filters = ["random(5)"]
```

**Parameters:**
- `N` (optional): Number of random entries to keep (default: 1)

---

### `reverse`
**Description:** Reverses the order of entries in the source.

**Usage:**
```toml
filters = ["reverse"]
```

**Parameters:** None

---

### `today`
**Description:** Keeps only entries that were published today.

**Usage:**
```toml
filters = ["today"]
```

**Parameters:** None

---

### `lasthours`
**Description:** Keeps only entries that are from the last X hours. Default is 24 hours.

**Usage:**
```toml
# Keep entries from last 24 hours (default)
filters = ["lasthours"]

# Keep entries from last 12 hours
filters = ["lasthours(12)"]
```

**Parameters:**
- `hours` (optional): Number of hours to look back (default: 24)

---

## Content Manipulation Filters

### `digest`
**Description:** Creates a digest of all entries, combining them into a single entry with proper heading structure.

**Usage:**
```toml
# Create digest with default heading level (2)
filters = ["digest"]

# Create digest with heading level 3
filters = ["digest(3)"]
```

**Parameters:**
- `level` (optional): Heading level to use (default: 2)

**Behavior:**
- Creates a single entry titled "Digest for [Source Title]"
- Each original entry becomes a section with its title as a heading
- Preserves content of all entries in order

---

### `combine`
**Description:** Similar to digest, but uses the first entry's title as the source title. Useful for merging sources.

**Usage:**
```toml
# Combine entries with default heading level (2)
filters = ["combine"]

# Combine entries with heading level 3
filters = ["combine(3)"]
```

**Parameters:**
- `level` (optional): Heading level to use (default: 2)

**Behavior:**
- Uses the first entry's title as the combined entry title
- Combines all entries into one
- Useful when merging multiple sources

---

### `links`
**Description:** Rewrites relative links to have proper protocol prefixes.

**Usage:**
```toml
filters = ["links"]
```

**Parameters:** None

**Behavior:**
- Converts `src="//example.com/image.jpg"` to `src="https://example.com/image.jpg"`
- Converts `href="//example.com/page"` to `href="https://example.com/page"`
- Also handles relative paths by prepending the source base URL

---

### `embedimage`
**Description:** Embeds images from entry attachments at specified positions.

**Usage:**
```toml
# Embed image at top (default)
filters = ["embedimage"]

# Embed image at bottom
filters = ["embedimage(bottom)"]

# Embed image at left with link
filters = ["embedimage(left,link)"]
```

**Parameters:**
- `position` (optional): Position to embed image - `top`, `bottom`, `left`, or `right` (default: `top`)
- `link` (optional): If "link" is specified, the image will be wrapped in a link to the entry

**Behavior:**
- Only works if the entry has an `ImageURL` attachment
- Adds appropriate CSS classes for positioning
- For left/right positions, adds clearfix for proper layout

---

### `replace`
**Description:** Replaces strings in entry content using configuration.

**Usage:**
```toml
filters = ["replace(myreplace)"]
```

**Configuration:**
```toml
[replace.myreplace]
from = "A string"
to = "Another string"
```

**Parameters:**
- `key`: The replacement configuration key to use

**Behavior:**
- Performs simple string replacement in all entry content
- Multiple replace configurations can be defined and used

---

### `includelink`
**Description:** Includes the link of entries in digest form.

**Usage:**
```toml
filters = ["includelink"]
```

**Parameters:** None

**Behavior:**
- When used with digest/combine filters, entry titles become clickable links
- Each entry title in the digest will link to the original entry URL

---

### `includesourcetitle`
**Description:** Includes source titles of entries in digest form.

**Usage:**
```toml
filters = ["includesourcetitle"]
```

**Parameters:** None

**Behavior:**
- When used with digest/combine filters, includes the source title above each entry
- Useful when merging multiple sources to maintain attribution

---

### `language`
**Description:** Keeps only entries in specified languages (best effort detection).

**Usage:**
```toml
# Keep only English entries
filters = ["language(en)"]

# Keep English and German entries
filters = ["language(en,de)"]
```

**Parameters:**
- `languages`: Comma-separated list of ISO 639-1 language codes

**Behavior:**
- Uses text analysis to detect language
- Only keeps entries that match the specified languages
- Works best with longer content for accurate detection

---

### `limitwords`
**Description:** Limits the number of words in each entry's content.

**Usage:**
```toml
filters = ["limitwords(100)"]
```

**Parameters:**
- `number`: Maximum number of words to keep

**Behavior:**
- Truncates content to the specified word count
- Attempts to truncate at sentence boundaries for better readability
- Preserves HTML structure while limiting text content

---

### `reskip`
**Description:** Skips entries whose titles match a regular expression.

**Usage:**
```toml
filters = ["reskip([Ss]ponsor.*)"]
```

**Parameters:**
- `regex`: Regular expression pattern to match against entry titles

**Behavior:**
- Removes entries whose titles match the provided regex
- Case-sensitive matching
- Logs debug information about skipped entries

---

## Specialized Filters

### `unseen`
**Description:** Keeps only unseen entries (entries not previously processed).

**Usage:**
```toml
filters = ["unseen"]
```

**Parameters:** None

**Behavior:**
- Tracks seen entries in a database (`goeland.db` by default)
- Only shows entries that haven't been seen before
- Marks entries as seen when they're processed
- Database location can be configured

---

### `retrieve`
**Description:** Retrieves full content from web pages using CSS selectors.

**Usage:**
```toml
# Retrieve content using CSS selector
filters = ["retrieve(div.content)"]
```

**Parameters:**
- `selector`: CSS selector to extract content from the page

**Behavior:**
- Fetches the original webpage for each entry
- Extracts content using the provided CSS selector
- Replaces entry content with the extracted content
- Fixes relative links and images to be absolute URLs
- Sanitizes HTML unless `--unsafe-no-sanitize-filter` is used

---

### `reddit`
**Description:** Provides better formatting for Reddit RSS feeds.

**Usage:**
```toml
filters = ["reddit"]
```

**Parameters:** None

**Behavior:**
- Improves formatting of Reddit posts
- Replaces thumbnail images with higher quality versions
- Handles both image posts and text posts
- Sanitizes HTML content for safety

---

### `untrack`
**Description:** Removes FeedBurner tracking pixels and links.

**Usage:**
```toml
filters = ["untrack"]
```

**Parameters:** None

**Behavior:**
- Removes FeedBurner tracking images & links

---

### `sanitize`
**Description:** Sanitizes HTML content of entries.

**Usage:**
```toml
filters = ["sanitize"]
```

**Parameters:** None

**Behavior:**
- Removes potentially dangerous HTML elements and attributes
- Should be used if `--unsafe-no-sanitize-filter` was passed to ensure safety of the relevant entries

---

### `toc`
**Description:** Creates a table of contents entry for all entries.

**Usage:**
```toml
# Create TOC with default title
filters = ["toc"]

# Create TOC using source title as link
filters = ["toc(title)"]
```

**Parameters:**
- `title` (optional): If "title" is specified, uses the source title as a link

**Behavior:**
- Creates a new entry at the beginning with a table of contents
- Lists all entry titles as clickable links
- Useful for long digests or combined sources

---

## Advanced Usage Examples

### `Creating a Daily Digest`
```toml
filters = ["today", "digest(3)", "includelink"]
```
- Keeps only today's entries
- Combines them into a digest with heading level 3
- Makes each entry title a clickable link

### `Language-Specific Newsletter`
```toml
filters = ["language(en)", "limitwords(150)", "embedimage(top)"]
```
- Keeps only English articles
- Limits each to 150 words
- Embeds images at the top of each entry

### `Reddit Image Gallery`
```toml
filters = ["reddit", "embedimage(bottom,link)"]
```
- Formats Reddit posts properly
- Embeds images at bottom with links to original posts

### `Multi-Source Merge`
```toml
filters = ["combine", "includesourcetitle", "toc(title)"]
```
- Combines multiple sources into one
- Shows source titles for attribution
- Adds table of contents at the top

## Filter Order Matters

The order of filters in your configuration is important:

1. **Filtering filters** (like `today`, `language`, `unseen`) should typically come first to reduce the number of entries early
2. **Content modification filters** (like `retrieve`, `replace`) should come next to work with individual entries
3. **Combining filters** (like `digest`, `combine`) should come last to work with the final set of entries

Example of good ordering:
```toml
filters = ["today", "language(en)", "retrieve(div.article)", "limitwords(200)", "digest(2)", "includelink"]
```
