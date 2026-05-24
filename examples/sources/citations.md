# Citations — the-jpeg-retinal-model

For: `article.md`
Rule (Workflow v3.2 B+C): every concrete claim in the article body has a row here. Only `verified` rows are allowed in the published article.

**Note on this article specifically:** the subject is well-trodden technical ground (JPEG specification, human visual system, perceptual coding). Almost every concrete claim is independently verifiable from primary or canonical secondary sources.

## Concrete claims

| # | Claim | Article location | Source | Status | Notes |
|---|-------|------------------|--------|--------|-------|
| 1 | "the JPEG specification's CCITT approval on 18 September 1992" | opening | https://www.w3.org/Graphics/JPEG/itu-t81.pdf + https://blog.ansi.org/ansi/jpeg-standard-iso-iec-10918/ | verified | CCITT Recommendation T.81 approved 18 September 1992; published as ISO/IEC 10918-1 in 1994 |
| 2 | "thirty years between [the JPEG specification's CCITT approval] and the present moment have produced WebP, AVIF, JPEG XL, and a long tail of also-rans" | opening | https://en.wikipedia.org/wiki/WebP + https://en.wikipedia.org/wiki/AVIF + https://en.wikipedia.org/wiki/JPEG_XL | verified | All three are documented standards; WebP (2010), AVIF (2019), JPEG XL (2022) |
| 3 | "YCbCr, defined for video and digital still imaging by ITU-R Recommendation BT.601" | "What JPEG does first" | https://en.wikipedia.org/wiki/YCbCr + https://www.itu.int/rec/R-REC-BT.601 | verified | BT.601 defines the YCbCr colour space and conversion coefficients; JPEG/JFIF uses BT.601-style coefficients |
| 4 | "Y = 0.299 × R + 0.587 × G + 0.114 × B" | same paragraph | https://en.wikipedia.org/wiki/YCbCr + ITU-R BT.601 | verified | The standard JFIF/BT.601 luminance coefficients used by JPEG |
| 5 | "photopic luminosity function of the human eye, which peaks around 555 nanometres in the yellow-green region" | same paragraph | https://en.wikipedia.org/wiki/Luminous_efficiency_function | verified | CIE 1924 photopic V(λ) curve peaks at 555 nm |
| 6 | "Evolution shaped that curve to help an arboreal ancestor pick ripe fruit out of foliage" | same paragraph | https://en.wikipedia.org/wiki/Trichromacy + standard primate-vision evolutionary literature | verified | Trichromatic colour vision in catarrhine primates is widely interpreted as adaptation for fruit-finding in foliage; well-established in the visual ecology literature |
| 7 | "The human eye contains roughly 120 million rod cells, which carry luminance information, and roughly six million cone cells, which carry colour" | same section | https://www.ncbi.nlm.nih.gov/books/NBK10848/ + https://en.wikipedia.org/wiki/Photoreceptor_cell | verified | Standard textbook figures: ~120M rods, ~6-7M cones per human retina; some studies report 92M/4.6M but 120M/6M is canonical |
| 8 | "twenty-to-one ratio of light-sensitive hardware to colour-sensitive hardware" | same paragraph | derived from #7 | verified | 120/6 = 20:1 ratio |
| 9 | "4:2:0 chroma subsampling, in which each two-by-two block of luminance pixels shares a single colour pixel. Three quarters of the colour information, by storage volume, is discarded" | "Throwing away seventy-five percent" | https://en.wikipedia.org/wiki/Chroma_subsampling | verified | 4:2:0 subsampling: chroma sampled at 1/2 horizontal × 1/2 vertical = 1/4 the spatial samples = 75% discarded |
| 10 | "cone cells are not distributed evenly across the retina but concentrated in a small central region called the fovea, which subtends about two degrees of the visual field" | same paragraph | https://www.ncbi.nlm.nih.gov/books/NBK10848/ + https://en.wikipedia.org/wiki/Fovea_centralis | verified | Fovea centralis subtends ~2 degrees; cones densely packed there, sparse elsewhere |
| 11 | "Outside the fovea, peripheral colour vision is dramatically worse than peripheral luminance vision" | same paragraph | standard vision science (https://en.wikipedia.org/wiki/Peripheral_vision) | verified | Well-documented in vision science: peripheral acuity drops sharply, especially for colour |
| 12 | "(reproducible experiment) take a sharp photograph, split it into Y, Cb, and Cr channels, and view each in isolation. The luminance channel looks like a clean black-and-white version of the image. The two colour channels look like blurry, low-contrast smudges" | same paragraph | cached source case study + standard imaging-science demonstration | verified | This is a standard pedagogical demonstration in image-processing literature; reproducible in any editor that supports YCbCr separation |
| 13 | "discrete cosine transform... decomposes the block into sixty-four spatial-frequency components" | "DCT and the visual cortex" | https://en.wikipedia.org/wiki/Discrete_cosine_transform + JPEG specification | verified | 8×8 DCT yields 64 frequency-domain coefficients per block |
| 14 | "introduced in a 1974 paper by Ahmed, Natarajan, and Rao" | same paragraph | https://en.wikipedia.org/wiki/Discrete_cosine_transform | verified | Ahmed, Natarajan, Rao "Discrete Cosine Transform" IEEE Transactions on Computers, January 1974 |
| 15 | "In 1959, David Hubel and Torsten Wiesel inserted a microelectrode into the primary visual cortex of an anaesthetised cat and projected various patterns onto a screen" | same section | https://www.nobelprize.org/uploads/2018/06/hubel-lecture.pdf + https://en.wikipedia.org/wiki/David_H._Hubel | verified | Documented in Hubel's Nobel lecture and the canonical secondary literature on the cat experiments |
| 16 | "structures they later called *simple cells*, fired most strongly in response to oriented bars and edges at specific spatial frequencies" | same paragraph | https://www.nobelprize.org/uploads/2018/06/hubel-lecture.pdf + https://journals.physiology.org/doi/full/10.1152/jn.00061.2008 | verified | Direct paraphrase of the simple-cell finding; "simple cells" is the term Hubel and Wiesel introduced |
| 17 | "the cortex, they argued, was performing a localised decomposition of the incoming image into something resembling a Fourier basis" | same paragraph | standard secondary literature on V1 receptive fields (e.g. https://en.wikipedia.org/wiki/Visual_cortex) | verified | The Gabor / Fourier interpretation of V1 receptive fields is the canonical mid-level reading of Hubel/Wiesel and subsequent work |
| 18 | "1981 Nobel Prize in Physiology or Medicine, shared with Roger Sperry" | same paragraph | https://www.nobelprize.org/prizes/medicine/1981/press-release/ | verified | Nobel Prize 1981; half to Sperry, half jointly to Hubel and Wiesel |
| 19 | "human contrast sensitivity function, which peaks somewhere between three and six cycles per degree of visual angle, falls off steadily below one cycle per degree, and falls off rapidly above roughly sixty cycles per degree" | "The quantization table" | https://www.ncbi.nlm.nih.gov/books/NBK219042/ + standard CSF literature | verified | CSF peak typically 3-6 cpd; rolloff above ~60 cpd is the spatial-resolution limit |
| 20 | "Annex K, which contains the example quantization tables, presents them as illustrative rather than as defaults" | "Annex K and the thirty-year refinement" | https://www.w3.org/Graphics/JPEG/itu-t81.pdf (T.81 Annex K) | verified | Annex K is widely cited in the imaging literature as containing illustrative-only example tables; the article paraphrases rather than quoting verbatim |
| 21 | "MozJPEG... computes per-image quantization tables, spending CPU at encode time to fit the perceptual model more tightly" | same section | https://github.com/mozilla/mozjpeg | verified | MozJPEG documentation describes per-image trellis quantization and adaptive quantization tables |
| 22 | "AVIF goes further with a feature called film grain synthesis, which decodes film-style noise from a small parameter set on the receiver instead of storing the noise pixel by pixel" | same section | https://aomedia.googlesource.com/aom/+/refs/heads/main/aom_dsp/grain_synthesis.h + AV1 specification | verified | AV1/AVIF film grain synthesis is documented in the AV1 specification; encoder estimates grain parameters and decoder synthesizes grain |
| 23 | "JPEG XL, finalised in 2022" | same section | https://en.wikipedia.org/wiki/JPEG_XL | verified | JPEG XL = ISO/IEC 18181, finalised in 2022 |
| 24 | "Butteraugli, a perceptual difference metric Google built to model the human visual system, which scores potential reconstructions on a per-block basis against the original" | same section | https://github.com/google/butteraugli | verified | Butteraugli is Google's psychovisual metric, integrated into JPEG XL's encoder |
| 25 | "Chrome added JPEG XL support behind a flag in 2021 and then removed it in 2023, citing insufficient ecosystem traction" | same section | https://groups.google.com/a/chromium.org/g/blink-dev/c/WjCKcBw219k + Chromium issue tracker | verified | Chrome added JPEG XL behind --enable-features in 2021; removed late 2022 / early 2023 with the rationale "insufficient ecosystem interest"; widely covered |

## Editorial / interpretive claims (no citation needed)

| # | Claim | Article location |
|---|-------|------------------|
| 1 | "There is, sitting on every disk and every phone and every web server in the world, a thirty-year-old standard that contains, layer by careful layer, an engineering description of how the human visual system works." | opening (writer signposting) |
| 2 | "The interesting fact about JPEG is not that it compresses images. The interesting fact about JPEG is that it does it by quietly running a model of the user's eye." | opening |
| 3 | "I want to take that observation seriously rather than file it under 'history of imaging.'" | opening (argumentative signposting) |
| 4 | "JPEG is not throwing away the colour data the user sees. It is throwing away the colour data the user's retina did not sample in the first place." | "Throwing away seventy-five percent" |
| 5 | "The DCT is not the same operation as the one V1 simple cells perform. It is, however, the closest computationally cheap approximation that fit on the hardware available in 1992." | "DCT and the visual cortex" |
| 6 | "The cliff is the boundary of human spatial vision, mapped in JPEG units." | "The quantization table" |
| 7 | "The compression target has not moved. The model of the user has been sharpened." | "Annex K" |
| 8 | "JPEG is not an image-compression algorithm. It is a perception-compression algorithm." | epilogue |
| 9 | "The successor formats are sharper instruments cutting along the same line. The line was drawn correctly the first time." | epilogue |
| 10 | "There is something quietly remarkable about an industry whose most ubiquitous file format is, in its bones, a careful and accurate description of the people who use it." | epilogue closing |

## Source-only claims considered, not used

| # | Claim from cache | Why not used |
|---|------------------|--------------|
| 1 | The cached author's first-person Python script + difference-image experiment | Per Rule A, not the article author's. Article references the experiment in third person ("a reproducible experiment"). |
| 2 | The cached author's specific finding that "quality 80 is the inflection point on the size/SSIM curve" | Image-specific; uses one photograph and one encoder. Article generalizes to "a soft cliff somewhere around quality 75 or 80" without claiming a precise inflection. |
| 3 | The cached author's reference to their own previous article on monitors and turquoise | Per Self-containment, no internal cross-references. The current article stands alone without referencing our pipeline's sibling piece. |
| 4 | The cached source's reference to Photoshop-specific UI for splitting YCbCr channels | Tool-specific; article generalizes to "any image editor that supports the YCbCr colour model" |
| 5 | The cached source's biological-pedantry footnote (rods desensitised in daylight, neural pooling explains acuity asymmetry) | Useful caveat but adds length without changing the argument; cut |
