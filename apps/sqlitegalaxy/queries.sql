-- Top 20 visited systems with 1st and last visit date
select count(*) as n, s.name, min(arrive), max(arrive)
from visits v join systems s on (v.sys = s.id)
where v.cmdr = 2
group by sys
order by n desc
limit 20
;

-- Which station did I visit
select p.name, s.name, count(*), min(d.arrive), max(d.arrive) as n
from docked d
join ports p on (d.port = p.id)
join systems s on (p.sys = s.id)
where d.cmdr = 1 and p.name like 'Fisher%'
group by port
order by n desc
;
