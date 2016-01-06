var TODAYLIST = (function(){
	'use strict'
	var that = {};
	that.gVar = {};
	that.init = function(){
		that.addpost.init();
	}
	return that;
})();

$(document).ready(function(){
	'use strict';
	TODAYLIST.init();
})