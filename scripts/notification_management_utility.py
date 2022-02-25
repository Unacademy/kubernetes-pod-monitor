import pymysql.cursors
import tabulate

host = input("Enter the MySQL hostname: ").strip()
user = input("Enter username: ").strip()
password = input("Enter Password: ").strip()
database = input("Enter database name (default kubernetes_pod_monitor): ").strip()

if database == "":
    database = "kubernetes_pod_monitor"


def print_table(table):
    if len(table) == 0:
        print("No results!")
    header = table[0].keys()
    rows = [x.values() for x in table]
    print(tabulate.tabulate(rows, header, tablefmt='grid'))


while True:
    print("\n\nChoose of the following options (enter 0 to exit):")
    print("1. View all slack channels configured.")
    print("2. View all ignored containers.")
    print("3. Turn off notifications for a container.")
    print("4. Configure slack channel for a namespace.")

    choice = input("Enter choice: ")
    try:
        choice = int(choice)
    except Exception as e:
        choice = -1
    if choice == 0:
        break
    if choice < 1 or choice > 4:
        print("Please enter a valid number!")
        continue

    connection = pymysql.connect(
        host=host,
        user=user,
        password=password,
        database=database,
        cursorclass=pymysql.cursors.DictCursor
    )
    with connection:
        with connection.cursor() as cursor:
            if choice == 1:
                sql = "SELECT clustername AS Cluster, namespace as Namespace, slack_channel as 'Slack channel' FROM `k8s_pod_crash_notify`"
                cursor.execute(sql)
                result = cursor.fetchall()
                print_table(result)
            if choice == 2:
                sql = "SELECT clustername AS Cluster, namespace as Namespace, containername AS 'Container name' FROM `k8s_crash_ignore_notify`"
                cursor.execute(sql)
                result = cursor.fetchall()
                print_table(result)
            if choice == 3:
                clustername = input("Enter cluster name: ")
                namespace = input("Enter namespace: ")
                containername = input("Enter container name to ignore: ")
                sql = "INSERT INTO `k8s_crash_ignore_notify` (`clustername`, `namespace`, `containername`) VALUES (%s, %s, %s)"
                try:
                    cursor.execute(sql, (clustername, namespace, containername))
                    connection.commit()
                    print(f"Notifications turned off for {containername} container.")
                except pymysql.err.IntegrityError as e:
                    print("Entry already exists!")
            if choice == 4:
                clustername = input("Enter cluster name: ")
                namespace = input("Enter namespace: ")
                channel = input("Enter slack channel to send alerts: ")
                sql = "INSERT INTO `k8s_pod_crash_notify` (`clustername`, `namespace`, `slack_channel`) VALUES (%s, %s, %s)"
                try:
                    cursor.execute(sql, (clustername, namespace, channel))
                    connection.commit()
                    print("Alerts configured.")
                except pymysql.err.IntegrityError as e:
                    print("Alerts already configured for this namespace!")

print("Utility exited.")
