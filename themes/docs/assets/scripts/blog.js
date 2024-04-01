'use strict';
var Blog = (function () {
    var siteData;
    var getCurrentTitle = function (siteData, pageLocation) {
        var title;
        siteData.every(function (page) {
            if (page.location === pageLocation) {
                title = page.title
                return false
            }
            if (pageLocation.startsWith(page.location) && page.children) {
                page.children.every(function (childPage) {
                    if (childPage.location === pageLocation) {
                        title = childPage.title
                        return false
                    }
                    return true
                })
            }
            return true
        });
        return title;
    }

    var getDoc = function (url, callback) {
        var pageLocation;
        if (url.length > 2) {
            pageLocation = url.substring(url.indexOf("/", 2) + 1, url.length)
        }
        if (!siteData) {
            App.require.require(MdRestConfig.BasePath + "mdrest_sitemap.json", "json", function (siteMap, err) {
                if (err) {
                    callback(null, {status: 404, msg: "没找到该文章"});
                    return
                }
                siteData = siteMap
                var title;
                if (url === '/') {
                    title = siteMap[0].title;
                    App.require.require(MdRestConfig.BasePath + siteMap[0].location + ".json", "json", function (pageData, err) {
                        if (err) {
                            callback(null, {status: 500, msg: "网络错误，网页未找到，请稍后重试"});
                            return
                        }
                        callback({
                            page: pageData.html,
                            title: pageData.title,
                            location: siteMap[0].location
                        })
                    });
                } else {
                    title = getCurrentTitle(siteData, pageLocation)
                }
                callback({
                    siteMap: siteMap,
                    title: title,
                    location: pageLocation,
                })
            });
        } else {
            callback({
                title: getCurrentTitle(siteData, pageLocation),
                location: pageLocation,
            })
        }
        if (!pageLocation) {
            return
        }
        App.require.require(MdRestConfig.BasePath + pageLocation + ".json", "json", function (data, err) {
            if (err) {
                callback(null, {status: 404, msg: "没找到该文章"});
                return
            }
            if (!data) {
                return;
            }
            callback({
                page: data.html,
                title: data.title,
                location: pageLocation,
            })
        });
    };
    return {
        getDoc: getDoc
    };
})();

