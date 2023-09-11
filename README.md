# memoria
Simple y ligero almacenamiento caché en memoria seguro para llamadas concurrentes.

## Ejemplo de uso:

```go
import "github.com/saulortega/memoria"

func main() {
    // Crear un almacén cuyos objetos almacenados serán eliminados
    // al transcurrir 20 segundos desde la última vez que fueron adquiridos
    // o al transcurrir cinco minutos desde que fueron almacenados, lo primero que ocurra.
    almacén := memoria.NuevoAlmacén(5*time.Minute, 20*time.Second)

    // Almacenar un objeto en memoria.
    almacén.Almacenar("identificador", "objeto")

    //...

    // Adquirir un objeto desde memoria.
    obj, existe := almacén.Adquirir("identificador")
    if !existe {
        // Manejar caso. Aquí obj es nil
        return
    }

    fmt.Println(obj.(string) == "objeto")
}
```
