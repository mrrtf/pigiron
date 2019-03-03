server = "localhost:8080"

function buildURL(what, query) {
        return "http://" + server + "/" + what + "?" + query
}

dualsampasJSON = ""
allds = ""
colorScheme = d3.scaleSequential(d3.interpolateYlGnBu)

function show(deid, bending) {

        var query = "deid=" + deid + "&bending=" + bending
        const dsURL = buildURL("dualsampas", query)
        const deURL = buildURL("degeo", query)

        let promiseArr = [d3.json(deURL), d3.json(dsURL)]

        Promise.all(promiseArr)
                .then(function(d) {
                        degeoJSON = d[0]
                        dualsampasJSON = d[1].DualSampas
                        changeValues()
                        createSVG(degeoJSON, dualsampasJSON)
                })
}

function changeValues() {

        for (i = 0; i < dualsampasJSON.length; i++) {
                dualsampasJSON[i].Value = Math.random()
        }

}

function updateValues() {
        
        allds.attr("fill", function(d) {
                return colorScheme(d.Value)
        })
}

function createSVG(degeoJSON, dualsampasJSON) {

        sx = degeoJSON.SX
        sy = degeoJSON.SY
        xleft = degeoJSON.X - sx / 2.0
        ybottom = degeoJSON.Y - sy / 2.0

        aspectRatio = 1.0 * sy / sx

        w = 800

        svg = d3.select("body").append("svg")
                .attr("width", w)
                .attr("height", (20 + aspectRatio * w))
                .attr("viewBox", "0 0 " + sx + " " + sy)

        dualsampas = svg.append("g").attr("class", "dualsampas")
                .attr("transform", "translate(" + (-xleft) + "," + (-ybottom) + ")")

        allds = dualsampas.selectAll("polygon").data(dualsampasJSON)
                .enter().append("polygon")
                .attr("stroke", "black")
                .attr("class", function(d) {
                        return "dualsampa DS" + d.ID
                })
                .attr("stroke-width", "0.1")
                // .attr("fill", function(d) {
                //         return colorScheme(d.Value)
                // })
                .on("mouseover", function(d) {
                        d3.select(this).attr("fill", "red")
                })
                .on("mouseout", function(d) {
                        d3.select(this).attr("fill", function(d) {
                                return colorScheme(d.Value)
                        })
                })
                .on("click", function(d) {
                        changeValues()
                        updateValues()
                })
                .attr("points", function(d) {
                        return d.Vertices.map(function(v) {
                                return [v.X, v.Y].join(",")
                        }).join(" ");
                })
        updateValues()
}
