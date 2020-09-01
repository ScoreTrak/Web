function numberRange (start, end) {
    if (start > end){
        const tmp = start
        start = end
        end = tmp
    }
    end = end+1
    return new Array(end - start).fill().map((d, i) => i + start);
}

function parse_index(rng){
    const stripped = rng.replace(/\s+/g, '')
    let ranges = stripped.split(',');
    let ret = []

    for (let i = 0; i < ranges.length; i++){
        if (ranges[i].split('-').length === 1){
            const num = parseInt(ranges[i])
            if (!isNaN(num)){
                ret.push(num)
            } else {
                return []
            }
        }
        else if (ranges[i].split('-').length === 2){
            const start = parseInt(ranges[i].split('-')[0])
            const end = parseInt(ranges[i].split('-')[1])
            if (!isNaN(start) && !isNaN(end)){
                ret.push(...numberRange(start, end))
            } else{
                return []
            }
        }
        else {
            return []
        }
    }
    return [...new Set(ret)].sort(function(a, b) {
        return a - b;
    })
}

export {parse_index, numberRange}