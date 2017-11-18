var map = null;
var map_maker = [];
var map_attr = '&copy; <a href="https://ditu.amap.com/">高德地图</a>'
var map_lng = 116.397697;
var map_lat = 39.906036;

var LeafIcon = L.Icon.extend({
    options: {
        shadowUrl: './img/shadow.png',
        iconSize: [25, 33],
        iconAnchor: [12, 33],
        popupAnchor: [1, -34],
        tooltipAnchor: [16, -28],
        shadowSize: [41, 41],
    }
});
var ofo_icon = new LeafIcon({iconUrl: './img/yellow.png'})
var mobike_icon = new LeafIcon({iconUrl: './img/red.png'})


function initMap() {
    map = L.map("map", {
        zoom: 15,
        center: [map_lat, map_lng],
        minZoom: 3,
        maxZoom: 18,
        // scrollWheelZoom: false,
    });


    map_layer = L.tileLayer.mapProvider('GaoDe.Normal.Map', {
        attribution: map_attr
    }).addTo(map)
}

function updateMap(ajax_data, new_lng, new_lat) {
    if (ajax_data == null) {
        return;
    }
    var ofo_list = ajax_data.ofo.bike_list;
    var mobike_list = ajax_data.mobike.bike_list;
    var total_list = [].concat(mobike_list).concat(ofo_list)

    if (map_maker.length !== 0) {
        for (var index in map_maker) {
            map.removeLayer(map_maker[index]);
        }
        map_maker = [];
    }

    for (var i in total_list) {
        var lat = total_list[i].lat;
        var lng = total_list[i].lng;
        if (total_list[i].car_type === "ofo") {
            var marker = L.marker([lat, lng], {icon: ofo_icon}).addTo(map);
        }
        else {
            var marker = L.marker([lat, lng], {icon: mobike_icon}).addTo(map);
        }
        map_maker.push(marker);
    };
    map.setView(new L.LatLng(new_lat, new_lng), 15);
}

function searchMap(addr_str) {
    var ajax_url = "/api/address/" + addr_str;
    var map_data = null;
    var bike_data = null;

    $.getJSON(ajax_url, function (data) {
        map_data = data;
    }).done(function () {
        if (map_data != null) {
            map_url = "/api/bike/" + map_data.lng + "/" + map_data.lat;
            $.getJSON(map_url, function (data) {
                bike_data = data;
            }).done(function () {
                updateMap(bike_data, map_data.lng, map_data.lat);
            });
        }
    }).fail(function () {
        showError();
    });
}

function addBike(lng, lat) {
    var ajax_url = "/api/bike/" + lng + "/" + lat;
    var ajax_data = null;
    $.getJSON(ajax_url, function (data) {
        updateMap(data, lng, lat);
    });
}

function initEvent() {
    $("#search").submit(function (event) {
        event.preventDefault();
        clearError();
        searchMap($('#bike').val())
    });

    $("#submit-btn").on("click", function(){
        console.log("niuniu");
        $("#search").submit();
    })
}

function showError() {
    $("#err-msg").show("fast");
    $("#bike").addClass("is-danger");
}

function clearError() {
    $("#err-msg").hide();
    $("#bike").removeClass("is-danger");
}

$(document).ready(function () {
    initMap();
    addBike(map_lng, map_lat);
    initEvent()
});
