#!/bin/bash
#
# Autor : Andreas Habegger <andreas.habegger@bfh.ch>
# Date  : 06-06-2014
# Desc  : Main build file
#
######################################################################

# Global variables
_MAINSHEET=${1##*/}
export MAINSHEET=${_MAINSHEET/.tex/}
export TEXINPUTS=$TEXINPUTS:../:../tplimages:../texmf:../pdf_tex:../exsrc:../pictures:../content:../database:../appendix
export WORKDIR="_sandbox"
export XFIGFILEDIR="fig"

# Local variables
#_LATEX=latex
_LATEX=pdflatex
_MAKEGLOSSARIES=makeglossaries
#_PDFVIEWER=acroread
_PDFVIEWER=evince
_VC_FILE="./tplscripts/vc"
_BUILDOUT="_output"
_LECTUREDIR="lectures"
_LECTUREMOD="lecture-"
_LECTURE=${2##*/}
_LECTURE=${_LECTURE/.tex/}
# List all required tools. System will check if the tools are installed
_REQ_SW="pdflatex gawk latex git"

function check_installation (){
  command -v $1 >/dev/null 2>&1 || return 1 && return 0 
}


function help {
  echo "First, input arg must be the lecture mode file lecture-*.tex file."
  echo "Second, the lecture to build."
  echo ""
  echo "-h    : shows this dialogue"
  echo "-l    : lists available lectures"
  echo "-m    : lists available modes"
  echo "-c    : convert xfig figures to eps"
  echo ""  
}

function lec_list {
  ls ${_LECTUREDIR}/lecture_*
  exit 0
}

function mod_list {
  ls ${_LECTUREMOD}*
  exit 0
}

function conv_xfig {
  cd ${XFIGFILEDIR}
  ls *.fig | while read filename; do fig2dev -L eps "$filename" "../${WORKDIR}/${filename%%'.fig'}.eps"; done
  cd ..
  exit 0
}

# Check if required tools are installed
for item in ${_REQ_SW[@]}; do
  check_installation $item
  if (( $? > 0 )) ; then
     echo >&2 "$item is not installed but required. Aborting...";
     exit 1
  fi
done 

# Menu dialogue
case "$1" in

-h)  help;
     exit 0
    ;;
-l)  lec_list
    ;;
-m)  mod_list
    ;;
-c)  conv_xfig
    ;;
*) echo "Signal number $1 is not processed"
   ;;
esac

if [ $# -lt 1 ]; then
        echo "Usage : $0 "
        echo ""
        help
        exit 0
fi

mkdir -p ${_BUILDOUT}

mkdir -p ${WORKDIR}


cd ${WORKDIR}

if [ -e ${MAINSHEET}.pdf ] ; then
 rm ${MAINSHEET}.pdf
fi

echo "\\def\\lectureToBuild{${_LECTURE}}" > env.tex

${_LATEX} ${MAINSHEET}.tex
cp -f "../database/*.bib" . && bibtex ${_MAINSHEET/.tex/.aux}
${_LATEX} ${MAINSHEET}.tex
${_MAKEGLOSSARIES} ${MAINSHEET}
${_LATEX} ${MAINSHEET}.tex
${_LATEX} ${MAINSHEET}.tex

if [ -e ${MAINSHEET}.pdf ] ; then
    cp ${MAINSHEET}.pdf ../${_BUILDOUT}/${MAINSHEET#*-}.pdf
else
   dvips $MAINSHEET.dvi
# 
#
# embedd all fonts in LaTeX to pdf as used in IEEE publications
# -------------------------------------------------------------
  ps2pdf14 -dPDFSETTINGS=/prepress -dEmbedAllFonts=true $MAINSHEET.ps && cp $MAINSHEET.pdf ../${_BUILDOUT}/${_LECTURE}-${MAINSHEET#*-}.pdf 
  ps2pdf $MAINSHEET.ps ../${_BUILDOUT}/${_LECTURE}-${MAINSHEET#*-}.pdf
#   ${_PDFVIEWER} ../${_BUILDOUT}/${_LECTURE}-${MAINSHEET#*-}.pdf
fi;

exit 0
