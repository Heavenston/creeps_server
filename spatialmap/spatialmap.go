package spatialmap

import "creeps.heav.fr/geom"

type Positioned interface {
    GetPosition() geom.Point
}

type SpatialMap struct {
    objects []Positioned
}

func (m *SpatialMap) Add(p Positioned) {
    m.objects = append(m.objects, p)
}

func (m *SpatialMap) Remove(p Positioned) {
    for i, o := range m.objects {
        if o == p {
            m.objects[i] = m.objects[len(m.objects)-1]
        }
        m.objects = m.objects[:len(m.objects)-1]
    }
}

func (m *SpatialMap) GetAt(point geom.Point) (Positioned, bool) {
    for _, obj := range m.objects {
        if obj.GetPosition() == point {
            return obj, true
        }
    }
    return nil, false
}

func (m *SpatialMap) GetIn(from geom.Point, upto geom.Point) (Positioned, bool) {
    for _, obj := range m.objects {
        if obj.GetPosition().IsWithing(from, upto) {
            return obj, true
        }
    }
    return nil, false
}
