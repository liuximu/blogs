/**
 * [main module]
 * @type {Object}
 */
const pageFooter = module.exports = {
	/** Map of new style */
	book: {
		assets: './assets',
		css: [
			'footer.css'
		]
	},

	/** Map of hooks */
	hooks: {
		'page:before': function (page) {
            if( page.content.indexOf("<!--")  != 0) {
                return page
            }

            i = page.content.indexOf("-->");
            if (i == -1) {
                return page
            }

            comment = page.content.substr(4, i-4);
            configOption = {};
            rows = comment.split('\n')
            for (ri=0; ri< rows.length; ri++) {
                row = rows[ri].trim()
                if(row == "") {
                    continue
                }

                mi = row.indexOf(":")
                if(mi == -1) {
                    continue 
                }

                configOption[row.substr(0,mi).trim()] = row.substr(mi+1).trim()
            }

            page.content = page.content.substr(i+3);

            source = ""
            if(configOption.source) {
                source = '<span style="float:left"> <a herf="' + configOption.source + '">译·原文</a></span>'
            }
			const htmlContents = '\n\n'+
            '<footer class="footer">' +
                '<div class="footer__container--normal" alt="">' +
                    '<div class="footer__description--normal">'+
                        '<p class="paragraph footer__author--normal" style="color: #000 !important;"></p>' +
                           '<p>' +
                                source +
                                '<span style="float: right">' +
                                   ' By ' + 
                                   configOption.author +
                                   ' <br>At ' + 
                                   configOption.date + 
                                   ' <br><a href="http://liuximu.com">熙穆·和顺</a> ' + 
                                '</span>' +
                           '</p>' +
                        '</p>' +
                    '</div>' +
                '</div>' +
            '</footer>';

			/** add contents to the original content */
			page.content = page.content + htmlContents;

			return page;
		}
	},

	/** Map of new blocks */
	blocks: {},

	/** Map of new filters */
	filters: { }
};
