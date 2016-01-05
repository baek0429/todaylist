var TODAYLIST = (function(){
	'use strict'
	var that = {};
	that.init = function(){
		that.column.init();
	}
	return that;
})();

$(document).ready(function(){
	'use strict';
	TODAYLIST.init();
})