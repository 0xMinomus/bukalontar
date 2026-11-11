# Script-scope evaluation inputs

These images are evaluation negatives, not landing-page samples or training examples. Input bytes are preserved so their SHA-256 can be matched to the JSONL runs; `/tmp/` paths in those logs refer to identical bytes.

| File | Source and attribution | License and changes |
| --- | --- | --- |
| `jawa-babad-surapati-page7.jpg` | [Babad Surapati, Acad. 240](https://commons.wikimedia.org/wiki/File:Babad_Surapati,_romantic_history_in_macapat_metres_Acad._240.pdf), Padang Kerta; Leiden University Libraries, [catalogue](http://hdl.handle.net/1887.1/item:2669790). Commons categorizes this as a manuscript in Javanese script; catalogue describes palm leaves and Javanese language. | [CC BY 4.0](https://creativecommons.org/licenses/by/4.0/). PDF page 7 rasterized with `pdftoppm -f 7 -singlefile -scale-to 1920 -jpeg`. Full page, no content crop. SHA-256: `67192017b51e9366ebdb08df02551f92e0628027a81df0874b4b304b47f48cbd`. |
| `kawi-sample.jpg` | [Kawi sample](https://commons.wikimedia.org/wiki/File:Kawi_sample.jpg), Aditya Bayu (Alteaven), 13 March 2014. Source describes Kawi script from Kakawin Ramayana 16, 31. | [CC BY-SA 3.0](https://creativecommons.org/licenses/by-sa/3.0/). Original JPEG, unchanged. SHA-256: `5989a89abbdb24b8a783264f7f07e58bb25f2def36fc8b6f95fa2dc7d2469069`. This is a modern typeset script sample with Latin labels, **not an ancient lontar photograph**. |

Positive controls use the two existing Bali-script photographs credited in [the app samples](../../../apps/web/public/samples/CREDITS.md). Script identity is based on source metadata and visual inspection, not an independent paleographer's certification. Kawi-script manuscript photographs have not been added to this set; the Kawi result must not be reported as a historical-lontar result.
