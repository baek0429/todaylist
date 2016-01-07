var TODAYLIST = (function(module){
	'use strict'
	module.addpost={};
	var that = module.addpost;
	that.init = function(){
		that.fileUpload();
	}
	that.fileUpload=function(){
		$("#input-id").fileinput();
	}
	return module;
})(TODAYLIST);