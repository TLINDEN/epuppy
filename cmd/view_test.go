package cmd

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPrepare(t *testing.T) {
	var tests = []struct {
		file string
		body string
	}{
		{
			"t/epub/basic-v3plus2.epub",
			`shirt court, an whinny retched a cordage offer groin-murder, picked inner windrow,`,
		},
		{
			"t/epub/childrens-literature.epub",
			`The child's natural literature. The world has lost certain secrets as the price of an advancing civilization. It is a commonplace of observation that no one can duplicate the success of Mother Goose, whether she be thought of as the maker of jingles or the teller of tales. The conditions of modern life preclude the generally naïve attitude that produced the folk rhymes, ballads, tales, proverbs, fables, and myths. The folk saw things simply and directly. The complex, analytic, questioning mind is not yet, either in or out of stories. The motives from which people act are to them plain and not mixed. Characters are good or bad. They feel no need of elaborately explaining their joys and sorrows. Such experiences come with the day's work. "To-morrow to fresh woods, and pastures new." The zest of life with them is emphatic. Their humor is fresh, unbounded, sincere; there is no trace of cynicism. In folk literature we do not feel the presence of a "writer" who is mightily concerned about maintaining his reputation for wisdom, originality, or style. Hence the freedom from any note of straining after effect, of artificiality. In the midst of a life limited to fundamental needs, their literature deals with fundamentals. On the whole, it was a literature for entertainment. A more learned upper class may have concerned itself then about "problems" and "purposes," as the whole world does now, but the literature of the folk had no such interests.`,
		},
		{
			"t/epub/cole-voyage-of-life.epub",
			`Thomas Cole is regarded as the founder of the Hudson River School, an American art movement that flourished in the mid-19th century and was concerned with the realistic and detailed portrayal of nature but with a strong influence from Romanticism. This group of American landscape painters worked between about 1825 and 1870 and shared a sense of national pride as well as an interest in celebrating the unique natural beauty found in the United States. The wild, untamed nature found in America was viewed as its special character; Europe had ancient ruins, but America had the uncharted wilderness. As Cole's friend William Cullen Bryant sermonized in verse, so Cole sermonized in paint. Both men saw nature as God's work and as a refuge from the ugly materialism of cities. Cole clearly intended the Voyage of Life to be a didactic, moralizing series of paintings using the landscape as an allegory for religious faith.`,
		},
		{
			"t/epub/epub30-spec.epub",
			`IDPF Members

Invited Experts/Observers

Version 2.0.1 of this specification was prepared by the International Digital Publishing Forum’s EPUB Maintenance Working Group under the leadership of:

Active members of the working group at the time of publication of revision 2.0.1 were:

Version 1.0 of this specification was prepared by the International Digital Publishing Forum’s Unified OEBPS Container Format Working Group under the leadership of:

Active members of the working group at the time of publication of revision 1.0 were:`,
		},
		{
			"t/epub/epub_sample_file_50KB.epub",
			`magna aliqua. Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna`,
		},
		{
			"t/epub/Fundamental-Accessibility-Tests-Basic-Functionality-v2.0.0.epub",
			`Open the Fundamental Accessibility Test book in the reading system.

If the test book is not available in the bookshelf, then open any other book that is available. 

If the reading system also supports side loading, then please provide notes about the accessibility of the side loading feature.

Indicate Pass or Fail.`,
		},
		{
			"t/epub/georgia-cfi.epub",
			`The Great Valley Region consists of folded sedimentary rocks, extensive erosion having removed the soft layers to form valleys, leaving the hard layers as ridges, both layers running in a N.E.-S.W. direction. In the extreme north-west corner of the state is a small part of the Cumberland Plateau, represented by Lookout and Sand Mts.

On the Blue Ridge escarpment near the N.E. corner of the state is a water-parting separating the waters which find their way respectively N.W. to the Tennessee river, S.W. to the Gulf of Mexico and S.E. to the Atlantic Ocean; indeed, according to B.M. and M.R. Hall (Water Resources of Georgia, p. 2), "there are three springs in north-east Georgia within a stone's throw of each other that send out their waters to Savannah, Ga., to Apalachicola, Fla., and to New Orleans, La." The water-parting between the waters flowing into the`,
		},
		{
			"t/epub/israelsailing.epub",
			` במשלוח דואר, מה שלפעמים היה נחמד כי גם חשבונות לתשלום לא היו מגיעים. 'טוב, מי`,
		},
		{
			"t/epub/jlreq-in-japanese.epub",
			` 2.5.1 基本版面からはみ出す例 2.5.2 基本版面で設定した行位置の適用 2.5.3 `,
		},
		{
			"t/epub/minimal-v2.epub",
			`This is a paragraph.`,
		},
		{
			"t/epub/minimal-v3.epub",
			`This is a paragraph.`,
		},
		{
			"t/epub/minimal-v3plus2.epub",
			`This is a paragraph.`,
		},
		{
			"t/epub/moby-dick.epub",
			`Call me Ishmael. Some years ago—never mind how long precisely—having little or no money in my purse, and nothing particular to interest me on shore, I thought I would sail about a little and see the watery part of the world. It is a way I have of driving off the spleen and regulating the circulation. Whenever I find myself growing grim about the mouth; whenever it is a damp, drizzly November in my soul; whenever I find myself involuntarily pausing before coffin warehouses, and bringing up the rear of every funeral I meet; and especially whenever my hypos get such an upper hand of me, that it requires a strong moral principle to prevent me from deliberately stepping into the street, and methodically knocking people’s hats off—then, I account it high time to get to sea as soon as I can. This is my substitute for pistol and ball. With a philosophical flourish Cato throws himself upon his sword; I quietly take to the ship. There is nothing surprising in this. If they but knew it, almost all men in their degree, some time or other, cherish very nearly the same feelings towards the ocean with me.`,
		},
		{
			"t/epub/sous-le-vent.epub",
			`SOUS LE VENT`,
		},
		{
			"t/epub/wasteland-otf.epub",
			`Line 20. Cf. Ezekiel 2:1.

23. Cf. Ecclesiastes 12:5.

31. V. Tristan und Isolde, i, verses 5-8.

42. Id. iii, verse 24.`,
		},
		{
			"pkg/epub/test.epub",
			`This EPUB file contains 10 hard coded page breaks. Note that these page breaks are different from the reflowed page numbers. If the total number of pages of this book in the reading app is not exactly 10, then you are looking at the reflowed pages. The app may be having a feature to switch between print page and reflowed pages for navigation. Or the print page list may be appearing in the TOC. If print page navigation feature does not exist or does not work then this test should be marked 'Fail'.`,
		},
	}

	for _, tt := range tests {
		testname := fmt.Sprintf("prepare/%s", tt.file)
		t.Run(testname, func(t *testing.T) {
			conf := Config{
				Document: "../" + tt.file,
			}

			ebook, err := Prepare(&conf)

			assert.NoError(t, err)
			assert.Contains(t, ebook.Body, tt.body, "expected text not found")
		})
	}
}
