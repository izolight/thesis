# Description - LaTex Course Build System

This is a very simple BFH LaTex beamer template.

##### Company
Bern University of Applied Sciences

##### Purpose
  This template is for creating BFH style LaTex beamer presentations.

##### Author
 * Andreas HABEGGER <andreas.habegger@bfh.ch>
 * Horst Heck <horst.heck@bfh.ch>


## Linux Prerequisites (Debian/Ubuntu)
Follow the procedure below to install 3rd party package (Debian/Ubuntu).

Install used LaTex packages:
```bash
apt-get install texlive-base texlive-extra-utils texlive-generic-recommended texlive-latex-base texlive-latex-extra
```

Install used fonts packages:
```bash
apt-get install texlive-fonts-extra texlive-fonts-recommended
```

## Usage
  - Set the variable values in variables.tex appropriate
  - Create content within **content** directory
  - Include the files in template_example.tex file
  - Rename template_example.tex file
  - Copy pictures / graphics to pictures directory. Include the graphics only by its appropriate name without 
 leading folder "pictures" and without file suffix.
  - Assemble the presentation either using your GUI LaTex IDE or on a Bash Shell by calling "build.do".

```bash
./build.do template_example.tex
```
  - To test build environment. Try to translate this example presentation. This will show missing packages or bad configuration in case of an error during build process.


##### Credits
Thanks to Horst Heck for an initial version of the LaTex BFH beamer theme.
