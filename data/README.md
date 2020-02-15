# Data Files

The Super Genki Database requires a few files for processing. These files are publically available with their own licences and because of that, not included in the source code.

## JMDict File
You'll need to download the dictionary file located at the [JMDict Homepage](http://edrdg.org/jmdict/j_jmdict.html). Super Genki supports the `JMDict` file, which supports English as well as the other language translations. Use the `gzip` utility to extract the XML file and copy it to the `data` directory.

```bash
$ gzip -d JMdict.gz
$ cp JMDict_e $GOPATH/src/github.com/royaloaklabs/super-genki-db/data/.
```

## Corpus Data Files
These are required for generating frequency data. Place all files in this directory.

### WorldLex
Download the Japanese Raw Frequencies files from http://worldlex.lexique.org/
```bash
$ unrar e Jap.Freq.2.rar
$ cp Jap.Freq.2.txt $GOPATH/src/github.com/royaloaklabs/super-genki-db/data/.
```

### Large Corpora used in CTS
Download the text files for **lemmas from the Internet corpus** and **word forms from the Internet corpus** from http://corpus.leeds.ac.uk/list.html

```bash
$ cp internet-jp-forms.num internet-jp.num $GOPATH/src/github.com/royaloaklabs/super-genki-db/data/.
```
