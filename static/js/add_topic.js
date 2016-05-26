/**
 * Created by Allen on 2016/5/21.
 */
$('.dropdown')
    .dropdown({
        // 你可以使用任何过渡
        transition: 'drop'
    })
;
$('.ui.accordion')
    .accordion()
;
var marked = require('marked');
marked.setOptions({
    renderer: new marked.Renderer(),
    gfm: true,
    tables: true,
    breaks: false,
    pedantic: false,
    sanitize: true,
    smartLists: true,
    smartypants: false
})
;

