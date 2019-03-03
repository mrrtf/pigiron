d3.selectAll('.polds')
        .on('mouseover', function(d) {
                d3.select(this)
                        .style('fill','red')
        })
        .on('mouseout', function(d) {
                d3.select(this)
                        .style('fill','none')
        })
                //
// d3.selectAll("polygon").on("mouseout", function() {
//         d3.select(this)
//                 .style("stroke-width","0.025px");
//         });
