# junit2html

Convert Junit XML reports (`junit.xml`) into HTML reports using a single standalone binary.

* Standalone binary.
* Failed tests are top, that's what's important.
* No JavaScript.
* Look gorgeous.

## Screenshot

![screenshot](screenshot.png)

## Usage

### Build

```bash
make build
```

### Run
```bash
./junit2html < <junit-file.xml> > <output-report.html>
```